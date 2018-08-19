## Peer 2 Peer Lending Application

A blockchain powered application that leverages smart contracts to transfer the money from lender to borrower automatically when certain criteria and conditions are met. These conditions are defined and executed as smart contractrts (chain code into Hyperledger).  Overall idea here is to remove or reduce role of intermediaries and make transactions smart and automatic.

For example, lender can offer $5000.00 amount for lending at 1% interest. Similarly, borrowr can also ask for borrowing money at specific interest rate. When lender and borrower's criteria are mactched with each other then money is automatically transfeered and transaction is recorded into the ledger. Here the idea is to remove the intermediaries between lender and borrower and make it automatic to reduce time and expenses. As borrower borrows money, its risk level is automatically increased as its liabilities are increased. Lender has to define condition to lend money with low or high risk level borrowres.

### Clarification:
This application and its entire idea is solely to learn developing Hyperledger Fabric based blockchain application. And the application is not really transferring any money :)

### Prerequisites and setup:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Node.js** v6.9.0 - 6.10.0 ( __Node v7+ is not supported__ )
* [Download Docker images](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#binaries)

### Application has 2 organisations, 1 Orderer, 2 CAs and 2 peers for each organisation.
* 2 CAs
* A SOLO orderer
* 4 peers (2 peers per Org)

### Application uses Node.js as web server which exposes REST APIs which invokes transactions on rquest from client application developed into Angular.
