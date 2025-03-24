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
3. Create a `.env` file and add your **Etherscan API Key**:  
   ```ini
   ETHERSCAN_API_KEY=your_api_key_here
   PORT=8080
   ```  
4. Run the API:  
   ```bash
   go run main.go
   ```  
5. Test the endpoints:  
   ```bash
   curl "http://localhost:8080/beneficiary?address=0xYourEthereumAddress"
   ```
