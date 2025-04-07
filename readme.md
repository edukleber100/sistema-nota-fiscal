# Microserviços - Impressão de Nota Fiscal

Este projeto simula o processo de cadastro de produtos, emissão e impressão de notas fiscais utilizando arquitetura de microserviços com Go (Golang).

## 🎯 Objetivo

Permitir:
- Cadastro de produtos com controle de saldo (estoque)
- Cadastro de notas fiscais com múltiplos produtos
- Impressão de notas fiscais (com validação de estoque e baixa)

## 🧱 Estrutura do Projeto

- **Serviço de Estoque** (`localhost:8081`)
  - `POST /produtos`: Cadastra um novo produto
  - `GET /produtos/listar`: Lista todos os produtos cadastrados
  - `POST /produtos/validar`: Valida se há saldo suficiente
  - `POST /produtos/baixar`: Dá baixa no estoque

- **Serviço de Faturamento** (`localhost:8082`)
  - `POST /notas`: Cadastra uma nova nota fiscal
  - `GET /notas/listar`: Lista as notas cadastradas
  - `GET /notas/imprimir/{id}`: Imprime uma nota (valida e baixa estoque)

## 🧪 Como Testar

1. **Inicie o serviço de estoque**
   ```bash
   go run estoque.go
   ```

2. **Inicie o serviço de faturamento**
   ```bash
   go run faturamento.go
   ```

3. **Cadastrando produtos**
   ```bash
   curl -X POST http://localhost:8081/produtos -H "Content-Type: application/json" -d '{"nome": "Produto A", "preco": 10.5, "saldo": 2}'
   ```

4. **Criando uma nota fiscal**
   ```bash
   curl -X POST http://localhost:8082/notas -H "Content-Type: application/json" -d '{"numero": "NF001", "produtos": [1]}'
   ```

5. **Imprimindo a nota fiscal**
   ```bash
   curl http://localhost:8082/notas/imprimir/1
   ```

## ✅ Comportamento Esperado

- Validação de estoque antes da impressão
- Atualização de saldo dos produtos
- Mudança do status da nota para "fechada"
- Feedback em caso de erro ou sucesso

## ⚙️ Tecnologias

- Golang
- HTTP (padrão REST)
- JSON (comunicação entre os serviços)

## 📌 Observações

- O sistema simula persistência em memória (sem banco de dados real).
- Os serviços são independentes, respeitando os princípios de microserviços.
- Caso ocorra falha na validação ou baixa de estoque, o usuário recebe um erro específico.
