use actix_web::{post, get, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

pub mod weathertweet {
    tonic::include_proto!("weathertweet");
}

use weathertweet::{
    weather_tweet_service_client::WeatherTweetServiceClient,
    WeatherTweetRequest, WeatherTweetResponse,
    Municipalities, Weathers,
};

#[derive(Debug, Serialize, Deserialize)]
struct WeatherInput {
    municipality: String,
    temperature: i32,
    humidity: i32,
    weather: String,
}

#[get("/")]
async fn root() -> impl Responder {
    HttpResponse::Ok().body("Bienvenido a la API REST para la app de clima")
}

#[post("/clima")]
async fn clima(data: web::Json<WeatherInput>) -> impl Responder {
    println!("Rust API recibió JSON: {:?}", data);

    let _ = dotenvy::dotenv().ok();

    let grpc_url = std::env::var("GRPC_SERVER_URL")
        .unwrap_or("http://0.0.0.0:50051".to_string());
    println!("Conectando a gRPC en: {}", grpc_url);

    let mut client = match WeatherTweetServiceClient::connect(grpc_url).await {
        Ok(c) => c,
        Err(e) => return HttpResponse::InternalServerError().body(format!("Error conectando: {}", e)),
    };

    // ✅ Convertir strings a enums generados por prost
    let municipality_enum = Municipalities::from_str_name(&data.municipality.to_lowercase());
    let weather_enum = Weathers::from_str_name(&data.weather.to_lowercase());

    if municipality_enum.is_none() || weather_enum.is_none() {
        return HttpResponse::BadRequest().body("Municipio o clima no válido según el proto");
    }

    let request = tonic::Request::new(WeatherTweetRequest {
        municipality: municipality_enum.unwrap() as i32,
        temperature: data.temperature,
        humidity: data.humidity,
        weather: weather_enum.unwrap() as i32,
    });

    match client.send_tweet(request).await {
        Ok(response) => {
            let resp: WeatherTweetResponse = response.into_inner();
            println!("Respuesta gRPC: {:?}", resp);
            HttpResponse::Ok().json(resp)
        }
        Err(e) => HttpResponse::InternalServerError().body(format!("Error gRPC: {}", e)),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenvy::dotenv().ok();

    match std::env::var("GRPC_SERVER_URL") {
        Ok(url) => println!("✅ Variable GRPC_SERVER_URL cargada: {}", url),
        Err(_) => println!("⚠️  GRPC_SERVER_URL no encontrada, usando valor por defecto"),
    }

    println!("🚀 Rust REST API corriendo en http://0.0.0.0:8080");
    HttpServer::new(|| App::new().service(root).service(clima))
        .bind(("0.0.0.0", 8080))?
        .run()
        .await
}
