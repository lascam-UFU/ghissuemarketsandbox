#!/bin/bash

# Environment variables
NETWORK_NAME="polar-n1"  # Update as needed
AUCTION_ID="repo123/issue123"
ISSUE_TITLE="Critical security vulnerability that needs to be patched."
AUCTION_METADATA="Required skills: Security, Go, Docker"
ESTIMATED_COST=100000  # in millisatoshis
BID_AMOUNTS=(95000 90000 85000)  # smaller bids in millisatoshis
BIDDERS=("bidder1" "bidder2" "bidder3")
MEMO="Payment for resolving security vulnerability"

# Get the OpenAI API Key from environment, using a default if not set
OPENAI_API_KEY=${OPENAI_API_KEY:-"your_default_openai_api_key_here"}
# Print the OpenAI API Key
echo "Using OpenAI API Key: $OPENAI_API_KEY"

# Containers
AGENTS=("polar-n1-agent1" "polar-n1-agent2" "polar-n1-agent3" "polar-n1-agent4")

# Function to check balances
check_balances() {
    echo "Checking balances for all agents..."
    for agent in "${AGENTS[@]}"; do
        echo "Balance for $agent:"
        docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" $agent ghissuemarket walletbalance
    done
}

# Print agent names and roles
echo "Agents involved in the experiment:"
echo "Auctioneer: ${AGENTS[0]}"
for (( i=1; i<${#AGENTS[@]}; i++ )); do
    echo "Bidder: ${AGENTS[i]}"
done

# Step 1: Check balances before the auction
check_balances

# Step 2: Add Issue
echo "Adding issue from ${AGENTS[0]}..."
docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket add-issue \
  --issue-id "issue123" \
  --issue-description "$ISSUE_TITLE" \
  --estimated-cost $ESTIMATED_COST \
  --metadata "$AUCTION_METADATA" || { echo "Failed to add issue"; exit 1; }

# Step 3: Outsource issue for bidding
echo "Outsourcing issue for bidding from ${AGENTS[0]}..."
docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket open-auction \
  --auction-id $AUCTION_ID \
  --issue-id "issue123" \
  --issue "$ISSUE_TITLE" \
  --starting-price $ESTIMATED_COST \
  --open-time $(date +%s) \
  --close-time $(( $(date +%s) + 300 )) \
  --metadata "$AUCTION_METADATA" || { echo "Failed to outsource issue"; exit 1; }

# Step 4: Place bids from bidders
echo "Placing bids from agents..."
for i in "${!BIDDERS[@]}"; do
    echo "Bidder ${BIDDERS[$i]} is placing a bid of ${BID_AMOUNTS[$i]} msat..."
    docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[$((i + 1))]} ghissuemarket place-bid \
      --auction-id $AUCTION_ID \
      --bidder-id "${BIDDERS[$i]}" \
      --bid-amount "${BID_AMOUNTS[$i]}" \
      --metadata "Skills: Security, Go, Docker" || { echo "Failed to place bid"; exit 1; }
done

# Step 5: Wait for the auction to close
echo "Waiting for auction to close..."
sleep 3 # Wait for the auction to close (adjust time as needed)

# Step 6: Close the auction
echo "Closing the auction from ${AGENTS[0]}..."
docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket close-auction \
  --auction-id $AUCTION_ID || { echo "Failed to close auction"; exit 1; }

# Step 7: Announce the winner
WINNING_BIDDER="${BIDDERS[1]}" # Assume bidder 2 wins
echo "Announcing the winner from ${AGENTS[0]}..."
docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket announce-winner \
  --auction-id $AUCTION_ID \
  --bidder-id "$WINNING_BIDDER" || { echo "Failed to announce winner"; exit 1; }

# Step 8: Create an invoice for the winning bidder
echo "Creating invoice for winning bidder from ${AGENTS[0]}..."
docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket add-invoice \
  --auction-id $AUCTION_ID \
  --amount 100000 --memo "$MEMO" || { echo "Failed to create invoice"; exit 1; }

# Step 9: Query for the payment request
echo "Querying for the payment request..."
PAYMENT_REQUEST=$(docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket query \
    "what is the payment request for auction_id='$AUCTION_ID' and its bid winner pubkey??")

# Print the payment request for copying
echo "Payment Request:"
echo "$PAYMENT_REQUEST"

# Step 10: Pay the invoice using the retrieved payment request
echo "Paying invoice from ${AGENTS[0]}..."
docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket pay-invoice \
  --payment-request "$PAYMENT_REQUEST" \
  --bidder-pubkey "your_bidder_pubkey_here" || { echo "Failed to pay invoice"; exit 1; }

# Step 11: Check balances after the auction
echo "Checking balances after the auction..."
check_balances

# Step 12: Handle no auction winner
if [[ -z "$WINNING_BIDDER" ]]; then
    echo "No auction winner found. Adding and resolving issue."
    docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket add-issue \
        --issue-id "issue124" \
        --issue-description "Resolve issue due to no auction winner." \
        --estimated-cost $ESTIMATED_COST \
        --metadata "Auto-resolved issue." || { echo "Failed to add issue"; exit 1; }

    echo "Resolving the issue from ${AGENTS[0]}..."
    docker exec -it -e OPENAI_API_KEY="$OPENAI_API_KEY" ${AGENTS[0]} ghissuemarket resolve-issue \
        --issue-id "issue124" \
        --resolution-details "Resolved internally due to no valid auction winner." || { echo "Failed to resolve issue"; exit 1; }
fi

echo "Experiment completed successfully."
