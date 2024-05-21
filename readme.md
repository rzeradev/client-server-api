# Desafio de Cotação do Dólar

Este projeto consiste em dois sistemas em Go (`client.go` e `server.go`) que interagem para obter a cotação do dólar, salvar em um banco de dados e retornar ao cliente. O projeto usa contextos para gerenciar timeouts e log de erros.

## Estrutura do Projeto

```
/client-server-api
|-- /server
|   |-- /repository
|   |   |-- cotacao.go
|   |-- /db
|-- /client
|   |-- client.go
|-- README.md
```

## Requisitos

- Go 1.21.3
- SQLite3
- Gorm

## Instruções para Rodar o Projeto

### Passo 1: Configurar o Servidor

1. Clone o repositório:

```sh
git clone https://github.com/rzeradev/client-server-api/server
cd client-server-api/server
```

2. Instale as dependências necessárias:

```sh
go mod tidy
```

3. Execute o servidor:

```sh
go run server.go
```

O servidor estará rodando na porta 8080.

### Passo 2: Configurar o Cliente

1. Em outra aba do terminal, navegue até o diretório do projeto:

```sh
cd client-server-api/client
```

2. Instale as dependências necessárias:

```sh
go mod tidy
```

3. Execute o cliente:

```sh
go run client.go
```

O cliente irá fazer uma requisição ao servidor, obter a cotação do dólar e salvar em um arquivo `cotacao.txt`.

## Detalhes Técnicos

### server.go

- O servidor fornece um endpoint `/cotacao` que busca a cotação do dólar de uma API externa.
- Usa `context` para gerenciar timeouts:
  - 200ms para buscar a cotação.
  - 10ms para salvar a cotação no banco de dados SQLite.
- As cotações são salvas no banco de dados com timestamp.

### client.go

- O cliente faz uma requisição ao servidor e obtém a cotação do dólar.
- Usa `context` para gerenciar um timeout de 300ms para a resposta do servidor.
- Salva a cotação no arquivo `cotacao.txt`.

### repository/cotacao.go

- Contém funções para inicializar o banco de dados e salvar cotações usando GORM.
- Usa `context` para operações de banco de dados com timeout.

## Logs de Erros

Os logs de erro são gerados se qualquer operação exceder o tempo limite configurado e são salvos no arquivo `logs.txt`.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).
