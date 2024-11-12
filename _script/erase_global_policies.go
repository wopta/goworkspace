package _script

/*
func EraseGlobalPolicies() {
	var policies []models.Policy

	iter := lib.WhereFirestore(lib.PolicyCollection, "company", "==", "global")
	if iter == nil {
		log.Fatalf("no policies found")
	}

	policies = models.PolicyToListData(iter)

	deleteBatch := make(map[string]map[string]interface{})

	deleteBatch[lib.PolicyCollection] = make(map[string]interface{})
	deleteBatch[lib.TransactionsCollection] = make(map[string]interface{})

	for _, policy := range policies {
		p := policy
		if p.Uid == "" {
			continue
		}

		transactions := transaction.GetPolicyTransactions("", p.Uid)

		deleteBatch[lib.PolicyCollection][p.Uid] = p

		for _, trx := range transactions {
			t := trx
			deleteBatch[lib.TransactionsCollection][t.Uid] = t
		}
	}

	err := lib.DeleteBatchFirestoreErr(deleteBatch)
	if err != nil {
		log.Fatal(err)
	}

	policyWhereClause := fmt.Sprintf(" WHERE uid IN ('" + strings.Join(lib.GetMapKeys(deleteBatch[lib.PolicyCollection]), "', '") + "')")
	log.Printf("policy where clause: %s", policyWhereClause)
	err = lib.DeleteRowBigQuery(lib.WoptaDataset, lib.PolicyCollection, policyWhereClause)
	if err != nil {
		log.Fatal(err)
	}

	transactionWhereClause := fmt.Sprintf(" WHERE uid IN ('" + strings.Join(lib.GetMapKeys(deleteBatch[lib.TransactionsCollection]), "', '") + "')")
	log.Printf("transaction where clause: %s", transactionWhereClause)
	err = lib.DeleteRowBigQuery(lib.WoptaDataset, lib.TransactionsCollection, transactionWhereClause)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d policies - %d transactions deleted", len(deleteBatch[lib.PolicyCollection]), len(deleteBatch[lib.TransactionsCollection]))
}
*/
