package cmd

import (
	"fmt"
	"iter"

	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"github.com/spf13/viper"
	"spheric.cloud/xiter"
)

func getAccountsFromConfig() ([]gaccount.GAccount) {
	var accounts []gaccount.GAccount
	err := viper.UnmarshalKey(accountsConfigKey, &accounts)
	if (err != nil) {
		showError("failed to read accounts from config:", err)
	}

	for i := range accounts { 
        err := accounts[i].Init(serviceName)
		if (err != nil) {
			showError(fmt.Sprintf("failed to initialize account '%s':", accounts[i].Name), err)
		}
	}

	return accounts
}

func getAccountsIterFromConfig() (iter.Seq[gaccount.GAccount]) {
	tempAccounts := getAccountsFromConfig()
	accounts := xiter.OfSlice(tempAccounts)
	return accounts
}