# go-viem
golang版本的viem.sh库


graph LR
A[Business Logic] --> B(ETHClient Interface)
B --> C[Production Client]
B --> D[Mock Client]