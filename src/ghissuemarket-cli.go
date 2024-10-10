// ghissuemarket-cli.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Log file paths
var (
	logFile        = "/var/log/ghissuemarket/ghissuemarket.log"
	privateLogFile = "/var/log/ghissuemarket/private.log"
	sysLogFile     = "/var/log/ghissuemarket/sys.log"
)

// Fixed LND paths for cert and macaroon
const (
	lndTlsCertPath  = "/home/lnd/.lnd/tls.cert"                                  // Path to tls.cert
	lndMacaroonPath = "/home/lnd/.lnd/data/chain/bitcoin/regtest/admin.macaroon" // Path to macaroon
)

// Auction structure
type Auction struct {
	UUID             string  `json:"uuid"` // UUID for grouping events
	AuctionID        string  `json:"auction_id"`
	IssueID          string  `json:"issue_id"` // Added field for issue ID
	Issue            string  `json:"issue"`
	StartingPrice    float64 `json:"starting_price"`
	OpenTime         int64   `json:"open_time"`
	CloseTime        int64   `json:"close_time"`
	AnnouncementTime int64   `json:"announcement_time"`
	Metadata         string  `json:"metadata"`
	State            string  `json:"state"`
	Timestamp        int64   `json:"timestamp"`
	AuctioneerPubKey string  `json:"auctioneer_pubkey"`
}

// Bid structure
type Bid struct {
	UUID         string  `json:"uuid"` // UUID for grouping events
	AuctionID    string  `json:"auction_id"`
	BidderID     string  `json:"bidder_id"`
	Amount       float64 `json:"amount"`
	Metadata     string  `json:"metadata"`
	State        string  `json:"state"`
	Timestamp    int64   `json:"timestamp"`
	BidderPubKey string  `json:"bidder_pubkey"`
}

// Invoice structure
type Invoice struct {
	UUID      string `json:"uuid"` // UUID for the invoice
	InvoiceID string `json:"invoice_id"`
	AuctionID string `json:"auction_id"`
	Amount    int64  `json:"amount"` // in satoshis
	Memo      string `json:"memo"`
	State     string `json:"state"`
	Timestamp int64  `json:"timestamp"`
}

// Issue structure
type Issue struct {
	UUID             string  `json:"uuid"` // UUID for the issue
	IssueID          string  `json:"issue_id"`
	IssueDescription string  `json:"issue_description"`
	EstimatedCost    float64 `json:"estimated_cost"`
	Metadata         string  `json:"metadata"`
	Timestamp        int64   `json:"timestamp"`
	State            string  `json:"state"` // e.g., Open, Resolved
}

// Main function
func main() {
	var rootCmd = &cobra.Command{
		Use:   "ghissuemarket",
		Short: "ghissuemarket is a CLI for managing decentralized auctions, bids, and issues",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Error: Unknown command or invalid usage.")
			showUsage(cmd)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Register all commands
	rootCmd.AddCommand(openAuctionCmd, closeAuctionCmd, placeBidCmd, announceWinnerCmd, addInvoiceCmd, payInvoiceCmd, walletBalanceCmd, queryCmd, addIssueCmd, resolveIssueCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error: Unknown command or invalid usage. Refer to help below:")
		rootCmd.Help()
		os.Exit(1)
	}
}

// open-auction command
var openAuctionCmd = &cobra.Command{
	Use:   "open-auction",
	Short: "Auctioneer opens a new auction with specific open and close times",
	Run: func(cmd *cobra.Command, args []string) {
		auctionID, _ := cmd.Flags().GetString("auction-id")
		issueID, _ := cmd.Flags().GetString("issue-id") // Added field for issue ID
		issue, _ := cmd.Flags().GetString("issue")
		startingPrice, _ := cmd.Flags().GetFloat64("starting-price")
		openTime, _ := cmd.Flags().GetInt64("open-time")
		closeTime, _ := cmd.Flags().GetInt64("close-time")
		metadata, _ := cmd.Flags().GetString("metadata")

		announcementTime := time.Now().Unix()
		auctioneerPubKey, err := getPublicKey("open-auction")
		if err != nil {
			logError(fmt.Sprintf("Failed to retrieve auctioneer public key for open-auction: %v", err))
			return
		}

		auction := &Auction{
			UUID:             uuid.New().String(), // Generate UUID for the auction
			AuctionID:        auctionID,
			IssueID:          issueID,
			Issue:            issue,
			StartingPrice:    startingPrice,
			OpenTime:         openTime,
			CloseTime:        closeTime,
			AnnouncementTime: announcementTime,
			Metadata:         metadata,
			State:            "Open",
			Timestamp:        announcementTime,
			AuctioneerPubKey: auctioneerPubKey,
		}

		saveDataToLog(auction, "auction-opened", logFile)
		fmt.Println(toJSON(auction))
	},
}

// close-auction command
var closeAuctionCmd = &cobra.Command{
	Use:   "close-auction",
	Short: "Auctioneer closes the auction",
	Run: func(cmd *cobra.Command, args []string) {
		auctionID, _ := cmd.Flags().GetString("auction-id")

		auctioneerPubKey, err := getPublicKey("close-auction")
		if err != nil {
			logError(fmt.Sprintf("Failed to retrieve auctioneer public key for close-auction: %v", err))
			return
		}

		entry := map[string]interface{}{
			"type":              "auction-closed",
			"auction_id":        auctionID,
			"auctioneer_pubkey": auctioneerPubKey,
			"timestamp":         time.Now().Unix(),
		}
		saveEntryToLog(entry, logFile)
		fmt.Printf("Auction %s closed.\n", auctionID)
	},
}

// place-bid command
var placeBidCmd = &cobra.Command{
	Use:   "place-bid",
	Short: "Bidder places a bid on an auction",
	Run: func(cmd *cobra.Command, args []string) {
		auctionID, _ := cmd.Flags().GetString("auction-id")
		bidderID, _ := cmd.Flags().GetString("bidder-id")
		bidAmount, _ := cmd.Flags().GetFloat64("bid-amount")
		metadata, _ := cmd.Flags().GetString("metadata")

		bidderPubKey, err := getPublicKey("place-bid")
		if err != nil {
			logError(fmt.Sprintf("Failed to retrieve bidder public key for place-bid: %v", err))
			return
		}

		bid := &Bid{
			UUID:         uuid.New().String(), // Generate UUID for the bid
			AuctionID:    auctionID,
			BidderID:     bidderID,
			Amount:       bidAmount,
			Metadata:     metadata,
			State:        "Placed",
			Timestamp:    time.Now().Unix(),
			BidderPubKey: bidderPubKey,
		}

		saveDataToLog(bid, "bid-placed", logFile)
		fmt.Println(toJSON(bid))
	},
}

// announce-winner command
var announceWinnerCmd = &cobra.Command{
	Use:   "announce-winner",
	Short: "Auctioneer announces the winner of an auction",
	Run: func(cmd *cobra.Command, args []string) {
		auctionID, _ := cmd.Flags().GetString("auction-id")
		bidderID, _ := cmd.Flags().GetString("bidder-id")

		auctioneerPubKey, err := getPublicKey("announce-winner")
		if err != nil {
			logError(fmt.Sprintf("Failed to retrieve auctioneer public key for announce-winner: %v", err))
			return
		}

		entry := map[string]interface{}{
			"type":              "winner-announced",
			"auction_id":        auctionID,
			"bidder_id":         bidderID,
			"auctioneer_pubkey": auctioneerPubKey,
			"timestamp":         time.Now().Unix(),
		}
		saveEntryToLog(entry, logFile)
		fmt.Printf("Winner announced for auction %s: Bidder %s\n", auctionID, bidderID)
	},
}

// add-invoice command
var addInvoiceCmd = &cobra.Command{
	Use:   "add-invoice",
	Short: "Bidder creates an invoice for the auctioneer to pay",
	Run: func(cmd *cobra.Command, args []string) {
		auctionID, _ := cmd.Flags().GetString("auction-id")
		amountFloat, _ := cmd.Flags().GetFloat64("amount")
		memo, _ := cmd.Flags().GetString("memo")

		// Convert amount to satoshis
		amount := int64(amountFloat * 1000) // Assuming amount is in millisatoshis

		lncliCmd := exec.Command("lncli",
			"--tlscertpath", lndTlsCertPath,
			"--macaroonpath", lndMacaroonPath,
			"addinvoice",
			fmt.Sprintf("%d", amount),
		)
		lncliOutput, err := lncliCmd.CombinedOutput()
		if err != nil {
			logError(fmt.Sprintf("Failed to create invoice using lncli: %v, Output: %s", err, string(lncliOutput)))
			return
		}

		var lncliResult map[string]interface{}
		if err := json.Unmarshal(lncliOutput, &lncliResult); err != nil {
			logError(fmt.Sprintf("Failed to parse lncli output: %v", err))
			return
		}

		invoiceID := fmt.Sprintf("invoice-%s-%s", auctionID, uuid.New().String())
		invoice := &Invoice{
			UUID:      uuid.New().String(), // Generate UUID for the invoice
			InvoiceID: invoiceID,
			AuctionID: auctionID,
			Amount:    amount,
			Memo:      memo,
			State:     "Created",
			Timestamp: time.Now().Unix(),
		}

		saveDataToLogWithMetadata(invoice, "invoice-created", lncliResult, logFile)
		fmt.Println(toJSON(invoice))
	},
}

// pay-invoice command
var payInvoiceCmd = &cobra.Command{
	Use:   "pay-invoice",
	Short: "Auctioneer pays an invoice, opens/closes channels",
	Run: func(cmd *cobra.Command, args []string) {
		paymentRequest, _ := cmd.Flags().GetString("payment-request")
		bidderPubKey, _ := cmd.Flags().GetString("bidder-pubkey")

		// Open channel
		openChannelCmd := exec.Command("lncli",
			"--tlscertpath", lndTlsCertPath,
			"--macaroonpath", lndMacaroonPath,
			"openchannel",
			bidderPubKey,
			"1000000000000", // Adjust the channel size as needed
		)
		openChannelOutput, err := openChannelCmd.CombinedOutput()
		if err != nil {
			logError(fmt.Sprintf("Failed to open channel: %v, Output: %s", err, string(openChannelOutput)))
			return
		}

		// Pay invoice
		lncliPayCmd := exec.Command("lncli",
			"--tlscertpath", lndTlsCertPath,
			"--macaroonpath", lndMacaroonPath,
			"payinvoice",
			paymentRequest,
		)
		lncliPayOutput, err := lncliPayCmd.CombinedOutput()
		if err != nil {
			logError(fmt.Sprintf("Failed to pay invoice: %v, Output: %s", err, string(lncliPayOutput)))
			return
		}

		// Confirm payment using lookupinvoice
		lookupInvoiceCmd := exec.Command("lncli",
			"--tlscertpath", lndTlsCertPath,
			"--macaroonpath", lndMacaroonPath,
			"lookupinvoice",
			paymentRequest,
		)
		lookupOutput, err := lookupInvoiceCmd.CombinedOutput()
		if err != nil {
			logError(fmt.Sprintf("Failed to lookup invoice: %v, Output: %s", err, string(lookupOutput)))
			return
		}

		// Parse lookup output to check payment status
		var lookupResult map[string]interface{}
		if err := json.Unmarshal(lookupOutput, &lookupResult); err != nil {
			logError(fmt.Sprintf("Failed to parse lookup invoice output: %v", err))
			return
		}

		// Check the payment status
		if paid, ok := lookupResult["state"].(string); ok && paid == "SETTLED" {
			fmt.Printf("Invoice paid successfully. Payment request: %s\n", paymentRequest)
		} else {
			fmt.Println("Payment not settled yet or failed.")
		}

		// Close channel
		closeChannelCmd := exec.Command("lncli",
			"--tlscertpath", lndTlsCertPath,
			"--macaroonpath", lndMacaroonPath,
			"closechannel",
			bidderPubKey,
		)
		closeChannelOutput, err := closeChannelCmd.CombinedOutput()
		if err != nil {
			logError(fmt.Sprintf("Failed to close channel: %v, Output: %s", err, string(closeChannelOutput)))
			return
		}

		// Log the transaction
		event := map[string]interface{}{
			"type":            "invoice-paid",
			"payment_request": paymentRequest,
			"metadata":        string(lncliPayOutput),
			"timestamp":       time.Now().Unix(),
		}
		saveEntryToLog(event, privateLogFile)

		// Automatically log wallet balance after transaction
		logWalletBalance()
	},
}

// walletbalance command
var walletBalanceCmd = &cobra.Command{
	Use:   "walletbalance",
	Short: "Check wallet balance and log to private.log",
	Run: func(cmd *cobra.Command, args []string) {
		balance := logWalletBalance()
		fmt.Printf("Wallet balance: %s\n", balance)
	},
}

// query command
var queryCmd = &cobra.Command{
	Use:   "query [query-string]",
	Short: "Query the environment for feedback",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		queryStr := args[0]
		fmt.Printf("Executing query: %s\n", queryStr)

		// Execute the ghissuemarket-feedback_engine command with the query
		out, err := exec.Command("/usr/local/bin/ghissuemarket-feedback_engine", queryStr).Output()

		if err != nil {
			logError(fmt.Sprintf("Error executing feedback engine: %v, Output: %s", err, string(out)))
			fmt.Println("Error executing query. Please check the logs for more details.")
			return
		}

		// If output is empty, notify the user
		if len(out) == 0 {
			fmt.Println("No response from the feedback engine.")
			return
		}

		// Display the query response
		fmt.Println("Query response:", string(out))
	},
}

// add-issue command
var addIssueCmd = &cobra.Command{
	Use:   "add-issue",
	Short: "Creates a new issue with an estimated cost",
	Run: func(cmd *cobra.Command, args []string) {
		issueID, _ := cmd.Flags().GetString("issue-id")
		issueDescription, _ := cmd.Flags().GetString("issue-description")
		estimatedCost, _ := cmd.Flags().GetFloat64("estimated-cost")
		metadata, _ := cmd.Flags().GetString("metadata")

		issue := &Issue{
			UUID:             uuid.New().String(), // Generate UUID for the issue
			IssueID:          issueID,
			IssueDescription: issueDescription,
			EstimatedCost:    estimatedCost,
			Metadata:         metadata,
			Timestamp:        time.Now().Unix(),
			State:            "Open",
		}

		saveDataToLog(issue, "issue-created", logFile)
		fmt.Println(toJSON(issue))
	},
}

// resolve-issue command
var resolveIssueCmd = &cobra.Command{
	Use:   "resolve-issue",
	Short: "Resolve an issue directly without outsourcing",
	Run: func(cmd *cobra.Command, args []string) {
		issueID, _ := cmd.Flags().GetString("issue-id")
		resolutionDetails, _ := cmd.Flags().GetString("resolution-details")

		event := map[string]interface{}{
			"type":               "issue-resolved",
			"issue_id":           issueID,
			"resolution_details": resolutionDetails,
			"timestamp":          time.Now().Unix(),
		}
		saveEntryToLog(event, logFile)
		fmt.Printf("Issue %s resolved with details: %s\n", issueID, resolutionDetails)
	},
}

// Utility functions

// saveDataToLog appends a JSON-encoded entry to the specified log file.
func saveDataToLog(data interface{}, eventType string, logFile string) {
	entry := map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      data,
	}
	saveEntryToLog(entry, logFile)
}

// saveDataToLogWithMetadata logs data with additional metadata.
func saveDataToLogWithMetadata(data interface{}, eventType string, metadata map[string]interface{}, logFile string) {
	entry := map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      data,
		"metadata":  metadata,
	}
	saveEntryToLog(entry, logFile)
}

// saveEntryToLog appends a JSON-encoded entry to the specified log file.
func saveEntryToLog(entry map[string]interface{}, logFile string) {
	logFileHandle, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", logFile, err)
	}
	defer logFileHandle.Close()

	eventJSON, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("Failed to marshal log entry: %v", err)
	}

	logFileHandle.Write(eventJSON)
	logFileHandle.Write([]byte("\n"))
}

// toJSON returns a pretty-printed JSON string of the given data.
func toJSON(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal to JSON: %v", err)
	}
	return string(jsonData)
}

// logError logs an error message to sys.log with a timestamp and unique event ID.
func logError(message string) {
	event := map[string]interface{}{
		"type":      "error",
		"timestamp": time.Now().Unix(),
		"message":   message,
	}
	saveEntryToLog(event, sysLogFile)
}

// getPublicKey retrieves the public key using lncli.
func getPublicKey(command string) (string, error) {
	cmd := exec.Command("lncli",
		"--tlscertpath", lndTlsCertPath,
		"--macaroonpath", lndMacaroonPath,
		"getinfo",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("lncli getinfo failed for %s: %v", command, err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("failed to parse lncli getinfo output for %s: %v", command, err)
	}

	pubKey, ok := result["identity_pubkey"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract public key from lncli getinfo response for %s", command)
	}

	return pubKey, nil
}

// logWalletBalance retrieves the wallet balance using lncli and logs it to private.log.
func logWalletBalance() string {
	lncliCmd := exec.Command("lncli",
		"--tlscertpath", lndTlsCertPath,
		"--macaroonpath", lndMacaroonPath,
		"walletbalance",
	)
	lncliOutput, err := lncliCmd.Output()
	if err != nil {
		logError(fmt.Sprintf("Failed to fetch wallet balance: %v", err))
		return ""
	}

	event := map[string]interface{}{
		"type":      "wallet-balance",
		"data":      string(lncliOutput),
		"timestamp": time.Now().Unix(),
	}
	saveEntryToLog(event, privateLogFile)

	return string(lncliOutput)
}

// showUsage displays the usage information for the given command.
func showUsage(cmd *cobra.Command) {
	cmd.Help()
}

// Initialize flags
func init() {
	openAuctionCmd.Flags().String("auction-id", "", "Unique ID for the auction")
	openAuctionCmd.Flags().String("issue-id", "", "Unique ID for the issue") // Added field for issue ID
	openAuctionCmd.Flags().String("issue", "", "Issue to be auctioned")
	openAuctionCmd.Flags().Float64("starting-price", 0, "Starting price of the auction")
	openAuctionCmd.Flags().Int64("open-time", 0, "Unix timestamp for when the auction opens")
	openAuctionCmd.Flags().Int64("close-time", 0, "Unix timestamp for when the auction closes")
	openAuctionCmd.Flags().String("metadata", "", "Additional information (e.g., required skills)")

	closeAuctionCmd.Flags().String("auction-id", "", "ID of the auction to be closed")

	placeBidCmd.Flags().String("auction-id", "", "ID of the auction to place a bid on")
	placeBidCmd.Flags().String("bidder-id", "", "Unique bidder ID")
	placeBidCmd.Flags().Float64("bid-amount", 0, "Bid amount (in satoshis)")
	placeBidCmd.Flags().String("metadata", "", "Additional information (e.g., skills)")

	announceWinnerCmd.Flags().String("auction-id", "", "ID of the auction for which to announce the winner")
	announceWinnerCmd.Flags().String("bidder-id", "", "ID of the winning bidder")

	addInvoiceCmd.Flags().String("auction-id", "", "ID of the auction for which to create an invoice")
	addInvoiceCmd.Flags().Float64("amount", 0, "Amount of the invoice (in millisatoshis)")
	addInvoiceCmd.Flags().String("memo", "", "Memo describing the invoice")

	payInvoiceCmd.Flags().String("payment-request", "", "Payment request for the invoice to be paid")
	payInvoiceCmd.Flags().String("bidder-pubkey", "", "Public key of the bidder")

	addIssueCmd.Flags().String("issue-id", "", "Unique ID for the issue")
	addIssueCmd.Flags().String("issue-description", "", "Description of the issue")
	addIssueCmd.Flags().Float64("estimated-cost", 0, "Estimated cost to resolve the issue")
	addIssueCmd.Flags().String("metadata", "", "Additional information about the issue")

	resolveIssueCmd.Flags().String("issue-id", "", "ID of the issue to resolve")
	resolveIssueCmd.Flags().String("resolution-details", "", "Details of how the issue was resolved")
}
