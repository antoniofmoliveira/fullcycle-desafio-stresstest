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

- apresenta resumo contendo para cada segundo a quantidade de requisições, a quantidade de erros recebidas do endpoint, o tempo médio das requisições naquele segundo, a quantiade de vezes que não foi possível conectar por limitação local.

- apresenta resumo por status

- apresenta resumo em varios percentis de duração

## para testar

- para executar um servidor que gera erros

    `make serverwitherror`

- para executar um servidor que não gera erros

    `make serverwithouterror`

- para testar o endpoint "GET" "/hello"

    `make testget`

- para testar o endpoint "POST" "/hello"

    `make testpost`
