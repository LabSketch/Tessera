# PROMPT MESTRE — PROJECTO TESSERA

## 1. CONTEXTO DO PROJECTO

**Nome:** Tessera  
**Descrição:** Fluxo usando máquinas de estado  
**Repositório:** https://github.com/LabSketch/Tessera  
**Linguagem principal:** Go (Golang)  
**Data base:** Outubro 2026  

### Problema a resolver

Existe uma biblioteca de documentos em **Markdown** contendo diagramas **Mermaid** (`stateDiagram-v2`) que descrevem máquinas de estados finitos. O projecto precisa de um motor que:

1. Leia esses arquivos Markdown;
2. Extraia e interprete os blocos Mermaid de diagramas de estado;
3. Execute **scripts Lua** vinculados às transições ou estados;
4. Simule o fluxo do diagrama seguindo eventos até chegar ao estado final.

---

## 2. ESPECIFICAÇÃO FUNCIONAL

### 2.1 Entrada

- Arquivos `.md` com a seguinte estrutura mínima:

  ````markdown
  ---
  diagram: Título do Diagrama
  ---
  
  ```mermaid
  stateDiagram-v2
    [*] --> EstadoInicial
    EstadoInicial --> EstadoFinal : evento_x {{ script: "action.lua" }}
  ```
  ````

  ```mermaid
  stateDiagram-v2
    [*] --> EstadoInicial
    EstadoInicial --> EstadoFinal : evento_x [script: "action.lua"]
    
  ```

### 2.2 Saída

* Log do fluxo executado (estados visitados, eventos disparados);
* Estado final alcançado;
* Valores retornados pelos scripts Lua (se houver);
* Status de sucesso/falha do fluxo.

### 2.3 Comportamentos de Máquina de Estados

O motor deve suportar, no mínimo:

Tipo | Característica | Suporte necessário
-- | -- | --
Finitos	| Estados + transições acionadas por eventos | ✅ Obrigatório
Hierárquicos | Estados compostos com subestados | ✅ Obrigatório
Guardas | Condições booleanas nas transições | ⚠️ Desejável
Ações entry / do / exit | Código Lua ligado a estados/transições | ✅ Obrigatório
Histórico | Retorno ao último subestado ativo | ❌ Fora do escopo v1
Submáquina | Reutilização de comportamento | ❌ Fora do escopo v1
Concursos/Ortogonais | Regiões paralelas | ❌ Fora do escopo v1

## 3. ARQUITETURA TÉCNICA

### 3.1 Módulos principais

```
tessera/
├── cmd/
│   └── tessera/
│       └── main.go              # CLI: entrada, execução, saída
├── pkg/
│   ├── parser/                  # Extrai blocos Mermaid do Markdown
│   │   └── markdown.go
│   ├── mermaid/                 # Parse e AST do stateDiagram-v2
│   │   ├── parser.go
│   │   ├── ast.go               # Estado, Transição, Guarda, Ação
│   │   └── lexer.go
│   ├── engine/                  # Runtime da máquina de estados
│   │   ├── machine.go           # Estado atual, transições, execução
│   │   └── runner.go            # Loop de eventos e validação
│   ├── lua/                     # Integração com interpreter Lua
│   │   ├── bridge.go            # Bind Go↔Lua
│   │   └── sandbox.go           # Segurança e restrições
│   └── logger/
│       └── log.go               # Logs estruturados
├── docs/                        # Documentação do projeto
├── tests/                       # Unitários e integração
└── go.mod / go.sum
```

### 3.2 Bibliotecas de dependência

Função | Sugestão | Comentário
-- | -- | --
Parser Lua | gopher-lua (Yuin/gopher-lua) | Interpreter Lua puro em Go, sem C
Parser Markdown | goldmark ou blackfriday | Extração de blocos code fences
Parser Mermaid | custom ou nomnoml adaptado | stateDiagram-v2 é limitado; parser custom recomendado
CLI | cobra | Estrutura modular de comandos
Log | zerolog ou log/slog | Logs estruturados JSON/texto

### 3.3 Modelos de dados internos

**Estado**:
```Go
type State struct {
  ID         string
  EntryScript  []string // Lua executado ao entrar
  DoScript     []string // Lua executado durante permanência
  ExitScript   []string // Lua executado ao sair
}
````

**Transição**:
```Go
type Transition struct {
  From    string
  To      string
  Event   string
  Guard   string // expressão Lua → bool
  Action  string // script Lua
}
````

**Máquina**:
```Go
type Machine struct {
  States        map[string]*State
  Transitions   []*Transition
  CurrentState  string
  History       map[string]string // estado anterior por estado composto
}
```

## 4. REQUISITOS NÃO FUNCIONAIS

Categoria | Especificação
-- | --
Performance	| < 100ms para iniciar e simular diagramas ≤ 50 estados
Segurança | Scripts Lua rodam em sandbox; sem acesso a filesystem/red
Portabilidade | Binários estáticos para Linux/macOS/Windows
Logging | Output claro de cada passo da execução
Testing | Cobertura ≥ 80% nos módulos críticos (parser, engine)
Documentação | README com exemplo de uso + docs/API

## 5.CASOS DE USO PRIMÁRIOS

1. Execução direta:
```sh
   tessera run diagrams/pedido.md --event sequence
```
2. Modo debug:
```sh
   tessera run diagrams/pedido.md --debug
```
3. Testes automatizados:
```sh
   tessera test tests/integration/*.yaml
```
4. Exportação de trajetória:
```sh
   tessera run diagrams/pedido.md --output trajectory.json
```

## 6. RESTRIÇÕES E DECISÕES DE ESCOPO

Item | Decisão | Justificativa
-- | -- | --
Lua 5.1 vs 5.3 | Lua 5.1 via gopher-lua | Compatibilidade mais ampla
Parser Mermaid externo | Não usar | stateDiagram-v2 insuficientemente documentado; custom parser mais controlável
Histórico | v2.0 | Complexidade alta; fora MVP
Estados ortogonais | v2.0 | Requer arquitetura paralela
Submáquinas | v2.0 | Necessita módulo de composição
Integração CLI | Cobra desde v1 | Escalabilidade de comandos

## 7. CRITÉRIOS DE ACEITAÇÃO (v1.0)

- ✅ Parser extrai stateDiagram-v2 de Markdown corretamente
- ✅ Motor executa sequência de eventos até estado final ou bloqueio
- ✅ Scripts Lua vinculados a transições são executados e podem alterar estado
- ✅ Guards em Lua permitem/bloqueiam transições conforme resultado booleano
- ✅ Logs mostram caminho completo percorrido (trajetória)
- ✅ Build binário único e executável em 3 SOs
- ✅ Testes unitários cobrem parsing e engine > 80%
- ✅ README com 1 exemplo funcional end-to-end

## 8. PRÓXIMOS PROMPTS DE IMPLEMENTAÇÃO DERIVADOS

Este Prompt Mestre deve gerar os seguintes prompts filhos:
1. Parser Mermaid em Go — lexer/parser para stateDiagram-v2 → AST
2. Motor de execução de máquina de estados — FSM core com suporte a scripts
3. Integração Lua (gopher-lua) — bridge Go ↔ Lua e sandbox seguro
4. CLI em Go (cobra) — comandos run, test, debug
5. Testes de integração — casos com diagramas em [Diagramas](docs/Diagrams.md)
