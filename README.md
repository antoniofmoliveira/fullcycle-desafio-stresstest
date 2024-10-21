# Desafio FullCycle Stress Tester

Desenvolver um sistema de linha de comando (CLI) em Go para executar testes de carga em serviços web, fornecendo insights críticos sobre o desempenho e a robustez desses serviços.

## Solução

Implementação de app para obter indicadores RED de um endpoint, conforme a definição:

- Rate: The number of requests or events processed per unit of time.

- Errors: The frequency and types of errors occurring during request processing.

- Duration: The time taken to process requests, including latency and throughput.


A app executa uma quantidade de requests:

- armazena em memória os resultados dos requests

- aceita flag para a quantidade de testes

- aceita flag para o endpoint a ser testado

- aceita flag para o intervalo entre cada request em microsegundos. isso é necessário para não esgotar os recursos da máquina executando a app

- usa pool de conexões e implementa repetição de request até 3 vezes com intervalo de 1 microsegundo. no caso de insucesso registra como erro de rede. 

- apresenta resumo no final contendo para cada segundo a quantidade de requisições, a quantidade de erros recebidas do endpoint, o tempo médio das requisições naquele segundo, a quantiade de vezes que não foi possível conectar por limitação local.

## Testes

uma execução normal contra um endpoint sem rate limiter que simula randomicamente um erro na resposta a cada 100 requisições e apresenta latencia randômica de até 100 microsegundos por requisição.

```bash
$ go run main.go --numtests 400000 
2024/10/21 07:00:04 Tests starting...
 Running  400000  tests with interval  1  microseconds and endpoint  http://localhost:8080/hello
Min/Seg       Rate          Erro        Tempo Medio     Erro Rede
4           50,435           507         381.942µs             0
5          101,561          1017        1.652311ms             0
6          116,525          1115         247.819µs             0
7          107,051          1110         270.546µs             0
8           24,428           257         707.244µs             0
2024/10/21 07:00:08 Tests finished.
```

aqui pode-se observar o mesmo endpoint com o rate limiter configurado para 10.000 requisições por segundo

```bash
$ go run main.go --numtests 400000 
2024/10/21 08:49:14 Tests starting...
 Running  400000  tests with interval  1  microseconds and endpoint  http://localhost:8080/hello
Min/Seg       Rate           Error        Avg Time       Net Error
4914        11,919           1,968       807.233µs               0
4915        97,725          87,772      15.82696ms               0
4916       112,907         102,969      4.373471ms               0
4917       116,336         106,396       256.732µs               0
4918        61,113          61,113        251.45µs               0
2024/10/21 08:49:18 Tests finished.
```

aqui pode-se observar o efeito do rate limiter configurado para 10.000 requisições a cada 5 segundos.

```bash
$ go run main.go --numtests 400000 
2024/10/21 08:40:08 Tests starting...
 Running  400000  tests with interval  1  microseconds and endpoint  http://localhost:8080/hello
Min/Seg       Rate           Error        Avg Time       Net Error
4008        73,785          63,849       750.032µs               0
4009       108,537         108,537      18.999011ms              0
4010        56,399          51,349      872.80571ms              0
4011        24,524          19,661      1.725696512s             0
4012        21,793          21,793      586.862088ms             0
4013        86,177          86,177      51.574782ms              0
4014        28,785          28,785      1.417574ms               0
2024/10/21 08:40:14 Tests finished.
```

