# temperature-system-otel

## Descrição
O projeto consiste em um sistema de medição de temperatura, onde é possível pesquisar o CEP de uma cidade e obter a temperatura atual da mesma.
Além de monitorar a aplicação com OpenTelemetry através do Zipkin.

## ENV
Na raíz do projeto deve ser criado um arquivo `.env` com a seguinte variável de ambiente:
```shell
KEY_API_WEATHERAPI=<TOKEN>
```
Obs.: O `<TOKEN>` pode ser obtido no site [WeatherAPI](https://www.weatherapi.com/).

## Execução
Para executar o projeto, é necessário ter o Docker instalado na máquina. Em seguida, execute o seguinte comando na raiz do projeto para buildar a imagem do Docker:
```shell
make build
```
Para executar o programa, utilize o seguinte comando:
```shell
make run
```

## Request
```curl
curl --location 'localhost:8081/temperature?cep=<CEP>' \
--header 'Content-Type: application/json'
```

## Response
```json
{
    "localidade": "São Paulo",
    "temp_c": 19.2,
    "temp_f": 66.56,
    "temp_k": 292.34
}
```

## IMPORTANTE
Não foi realizado o deploy na Cloud Run, pois eu já tinha gastado todo meu crédito.