📊 Google Ads MCC Metrics Extractor
Sistema desenvolvido por Gabriel Netto

⚡ Performance e Arquitetura da Solução

A ideia central deste projeto foi criar um script extremamente rápido, capaz de processar grandes volumes de contas Google Ads em poucos segundos.
Por isso, a linguagem escolhida foi Go (Golang), aproveitando ao máximo o poder das goroutines e do paralelismo nativo oferecido pela linguagem.

Em vez de percorrer as contas de forma sequencial, o sistema distribui o trabalho entre múltiplos núcleos, criando rotinas independentes para cada subconta da MCC.
Isso permite:

Executar até 100 contas simultaneamente (respeitando a disponibilidade de CPU)

Realizar consultas pesadas com latência extremamente baixa

Saturar corretamente o uso de threads, sem bloqueios desnecessários

Escalonar horizontalmente com precisão

No meu caso real, com 62 contas dentro da MCC, o script conseguiu:

🔥 Processar 62 contas e mais de 200 campanhas em apenas ~1.7 segundos (1.700 ms)

Esse ganho de performance torna o sistema ideal para:

Agências com grande volume de clientes

Rotinas de coleta de métricas em tempo real

Automação de dashboards atualizados

Sistemas que exigem respostas rápidas em larga escala

A arquitetura permite ainda novas otimizações, como pools de workers, limites dinâmicos por CPU, caching interno e pré-carregamento de estruturas de requisição.


Você pode clonar e aumentar a query de acordo com suas necessidades e padrão SQL do google ads API.


✔️ Métricas detalhadas por campanha
✔️ Métricas agregadas por conta
✔️ Percentual de otimização da conta
✔️ Dados prontos para relatórios ou dashboards

🚀 Funcionalidades Principais
🔹 1. Acesso Automático às Subcontas

O sistema se conecta à MCC principal e itera por todas as subcontas vinculadas, sem necessidade de intervenção manual.

🔹 2. Coleta de Métricas de Campanhas

Para cada campanha é possível obter por padrão no código (voce pode aumentar) sempre para o mês atual:

Impressões

Cliques

Conversões

Custo total


🔹 3. % de Otimização da Conta
Retorna a % de otimização dada pelo google para cada conta

Histórico de desempenho

Esse percentual é retornado para cada subconta.

🔹 4. Saída Padronizada

Todos os dados são retornados em um formato fácil de consumir, como:

JSON


Dashboard (caso integrado a algum BI)

🛠️ Tecnologias Utilizadas

Google Ads API

Linguagem: GOlang

Autenticação OAuth2
