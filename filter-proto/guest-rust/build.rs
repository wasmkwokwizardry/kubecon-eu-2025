use std::io::Result;
fn main() -> Result<()> {
    prost_build::compile_protos(
        &["protos/filter/filter.proto"],
        &["../host", "../host/protos/include"],
    )?;
    Ok(())
}
