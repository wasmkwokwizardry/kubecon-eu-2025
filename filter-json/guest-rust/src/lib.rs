use extism_pdk::*;
use regex::Regex;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::hash::Hash;

typify::import_types!("k8s_api_v1.json");

#[derive(FromBytes, Deserialize, Debug)]
#[encoding(Json)]
pub struct FilterInput {
    pub pod: IoK8sApiCoreV1Pod,
    pub node_info: NodeInfo,
}

#[derive(FromBytes, Deserialize, Debug)]
#[encoding(Json)]
pub struct NodeInfo {
    pub node: IoK8sApiCoreV1Node,
    pub image_states: Option<HashMap<String, ImageStateSummary>>,
}

#[derive(FromBytes, Deserialize, Debug)]
#[encoding(Json)]
pub struct ImageStateSummary {
    pub size: u64,
    pub num_nodes: u32,
}

#[derive(ToBytes, Serialize, PartialEq, Debug)]
#[encoding(Json)]
pub struct Status {
    pub code: StatusCode,
    pub reason: String,
}

#[repr(i32)]
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(into = "i32")]
pub enum StatusCode {
    Success = 0,
    Error = 1,
    Unschedulable = 2,
    UnschedulableAndUnresolvable = 3,
    Wait = 4,
    Skip = 5,
}

impl From<StatusCode> for i32 {
    fn from(code: StatusCode) -> Self {
        code as i32
    }
}

const REGEX_ANNOTATION_KEY: &str = "scheduler.wasmkwokwizardry.io/regex";

#[plugin_fn]
pub fn filter(input: FilterInput) -> FnResult<Status> {
    Ok(filter_impl(input))
}

fn filter_impl(input: FilterInput) -> Status {
    let pattern = match input
        .pod
        .metadata
        .as_ref()
        .and_then(|m| m.annotations.get(REGEX_ANNOTATION_KEY))
    {
        Some(p) => p,
        None => {
            return Status {
                code: StatusCode::Success,
                reason: "no regex annotation found".to_string(),
            };
        }
    };

    let regex = match Regex::new(pattern) {
        Ok(r) => r,
        Err(e) => {
            return Status {
                code: StatusCode::Error,
                reason: format!("failed to compile regex \"{}\": {}", pattern, e),
            };
        }
    };

    let node_name = match input
        .node_info
        .node
        .metadata
        .as_ref()
        .and_then(|m| m.name.as_ref())
    {
        Some(n) => n,
        None => {
            return Status {
                code: StatusCode::Error,
                reason: "node name not found".to_string(),
            };
        }
    };

    if regex.is_match(node_name) {
        Status {
            code: StatusCode::Success,
            reason: format!("node \"{}\" matches regex \"{}\"", node_name, pattern),
        }
    } else {
        Status {
            code: StatusCode::Unschedulable,
            reason: format!(
                "node \"{}\" does not match regex \"{}\"",
                node_name, pattern
            ),
        }
    }
}
