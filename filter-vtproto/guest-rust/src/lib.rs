use extism_pdk::*;
use prost::Message;
use regex::Regex;

pub mod k8s {
    pub mod io {
        pub mod api {
            pub mod core {
                pub mod v1 {
                    include!(concat!(env!("OUT_DIR"), "/k8s.io.api.core.v1.rs"));
                }
            }
        }
        pub mod apimachinery {
            pub mod pkg {
                pub mod api {
                    pub mod resource {
                        include!(concat!(
                            env!("OUT_DIR"),
                            "/k8s.io.apimachinery.pkg.api.resource.rs"
                        ));
                    }
                }
                pub mod apis {
                    pub mod meta {
                        pub mod v1 {
                            include!(concat!(
                                env!("OUT_DIR"),
                                "/k8s.io.apimachinery.pkg.apis.meta.v1.rs"
                            ));
                        }
                    }
                }
                pub mod runtime {
                    include!(concat!(
                        env!("OUT_DIR"),
                        "/k8s.io.apimachinery.pkg.runtime.rs"
                    ));
                }
                pub mod util {
                    pub mod intstr {
                        include!(concat!(
                            env!("OUT_DIR"),
                            "/k8s.io.apimachinery.pkg.util.intstr.rs"
                        ));
                    }
                }
            }
        }
    }
}

pub mod filter {
    include!(concat!(env!("OUT_DIR"), "/filter.rs"));
}

use filter::{FilterInput, Status, StatusCode};

const REGEX_ANNOTATION_KEY: &str = "scheduler.wasmkwokwizardry.io/regex";

#[plugin_fn]
pub fn filter(input: Vec<u8>) -> FnResult<Vec<u8>> {
    Ok(FilterInput::decode(&input[..]).and_then(|input| {
        let status = filter_impl(input);
        let mut buf = Vec::new();
        buf.reserve(status.encoded_len());
        status.encode(&mut buf).unwrap();
        Ok(buf)
    })?)
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
                code: StatusCode::Success as i32,
                reason: "no regex annotation found".to_string(),
            };
        }
    };

    let regex = match Regex::new(pattern) {
        Ok(r) => r,
        Err(e) => {
            return Status {
                code: StatusCode::Error as i32,
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
                code: StatusCode::Error as i32,
                reason: "node name not found".to_string(),
            };
        }
    };

    if regex.is_match(node_name) {
        Status {
            code: StatusCode::Success as i32,
            reason: format!("node \"{}\" matches regex \"{}\"", node_name, pattern),
        }
    } else {
        Status {
            code: StatusCode::Unschedulable as i32,
            reason: format!(
                "node \"{}\" does not match regex \"{}\"",
                node_name, pattern
            ),
        }
    }
}
