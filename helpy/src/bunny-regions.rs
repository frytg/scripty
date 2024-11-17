use reqwest;
use serde_json::Value;
use std::env;
use std::fs::File;
use std::io::Write;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    fetch_and_save_runtimes().await?;
    Ok(())
}

async fn fetch_and_save_runtimes() -> Result<(), Box<dyn std::error::Error>> {
    println!("Fetching runtimes data...");

    // Get the API token from the environment variable
    let api_token = env::var("BUNNYNET_API_KEY")
        .expect("BUNNYNET_API_KEY environment variable not set");

    // Create a new reqwest client
    let client = reqwest::Client::new();

    // Make the API request
    let response = client
        .get("https://api.bunny.net/region")
        .header("AccessKey", api_token)
        .send()
        .await?;

    // Check if the request was successful
    if !response.status().is_success() {
        return Err(format!("API request failed with status: {}", response.status()).into());
    }

    // Parse the JSON response
    let json: Value = response.json().await?;

    // Create the tmp directory if it doesn't exist
    std::fs::create_dir_all("tmp")?;

    // Write the JSON to a file
    let mut file = File::create("tmp/bunny-regions.json")?;
    file.write_all(serde_json::to_string_pretty(&json)?.as_bytes())?;

    println!("Runtimes data has been saved to bunny-regions.json");

    Ok(())
}
