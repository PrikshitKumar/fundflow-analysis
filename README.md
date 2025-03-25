# **Ethereum Fund Flow Analysis API** 🚀  

This repository contains a **Go-based API** using the **Gin framework** to analyze Ethereum fund flows. The API fetches transactions from the **Etherscan API** and determines the **beneficiary addresses** (final recipients of funds) in complex transaction chains.  

## **📌 Features**  
- 🏦 **Beneficiary Analysis (`/beneficiary`)** – Tracks fund outflows and identifies recipients.  
- 💰 **Payer Analysis (`/payer`)** *(Bonus)* – Identifies sources of incoming funds.  
- 🔗 **Supports all transaction types** – Normal, internal, and token transfers (ERC-20, ERC-721, ERC-1155).  
- ⚡ **Optimized with concurrency** – Fast API calls using Goroutines.  
- 🔍 **Structured JSON Response** – Provides a detailed breakdown of fund movements.  

## **📖 Setup & Usage**  
1. Clone the repository:  
   ```bash
   git clone https://github.com/PrikshitKumar/fundflow-analysis
   cd fundflow-api
   ```  
2. Install dependencies:  
   ```bash
   go mod tidy
   ```  
3. Generate the API Key from: `https://etherscan.io/apidashboard` 
4. Create a `.env` file and add your **Etherscan API Key**:  
   ```bash
   ETHERSCAN_API_KEY=your_api_key_here
   PORT=8080
   ```  
5. Run the API:  
   ```bash
   go run main.go
   ```  
6. Test the endpoints for Beneficiary and payer:  
   ```bash
   curl -X GET "http://localhost:8080/beneficiary?address=0x1ecD55bD5C5754d44b88937928Faf00C8BDc4Ae8" -H "Content-Type: application/json" | jq
   curl -X GET "http://localhost:8080/beneficiary?address=0x2e5eF37Ade8afb712B8Be858fEc7389Fe32857e2" -H "Content-Type: application/json" | jq
   curl -X GET "http://localhost:8080/payer?address=0x1218E12D77A8D1ad56Ec2f6d3d09A428cb7FDA7c" -H "Content-Type: application/json" | jq
   ```
