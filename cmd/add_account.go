/*
Copyright © 2024 Eugene Shtoka <eshtoka@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"github.com/EugeneShtoka/figoro/lib/gaseed"
	"github.com/EugeneShtoka/figoro/lib/gauth"
	"github.com/EugeneShtoka/figoro/lib/typedkeyring"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultPort int = 58080
	minPort = 1024
	maxPort = 65535
	credFile string
	accountsConfigKey = "accounts"
)

var addAccountCmd = &cobra.Command{
	Use:   "account [account name (at least 3 letters)]",
	Args:  cobra.ExactArgs(1),
	Short: "Add account",
	Long: "Add account. Requires account name to delete",
	Run: func(cmd *cobra.Command, args []string) {
		err := addAccount(args[0], &logger)
		if (err != nil) {
			showError(fmt.Sprintf("failed to add account '%s'", args[0]), err)
			cmd.Usage()
		}
	},
}

func addAccount(accName string, logger *zerolog.Logger) error { 
	minLength := 3
	err := validateString("account name", accName, &minLength, nil)
	if (err != nil) {
		return err
	}

	accounts := getAccountsFromConfig()
	if (slices.ContainsFunc(accounts, func(acc gaccount.GAccount) bool { return acc.Name == accName })) {
		return fmt.Errorf("account '%s' already exists in config", accName)
	}

	clientID, err := getStringProperty("clientID", 60, 100)
	if (err != nil) {
		return err
	}
	clientSecret, err := getStringProperty("clientSecret", 30, 40)
	if (err != nil) {
		return err
	}
	port, err := getPort()
	if (err != nil) {
		return err
	}

	seed, err := authorize(clientID, clientSecret, fmt.Sprintf("%d", port), logger)
	if (err == nil) {
		err = saveAccount(accName, seed)
	}

	return err
}


func validatePort(portString string) error {
	var port, err = strconv.Atoi(portString)
	if (err != nil) {
		return errors.New("invalid input, please enter an int")
	}else if (port < minPort) {
		return fmt.Errorf("port must be greater than %d", minPort)
	} else if (port > maxPort) {
		return fmt.Errorf("port must be less than %d", maxPort)
	}
	return nil
}

func getPort() (int, error) {
	var port = viper.GetInt("port")
	var err = validatePort(fmt.Sprintf("%d", port))

	if  (err == nil) {
		return port, nil
	}
	showError(fmt.Sprintf("invalid config for port: %d", port), err)

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Enter port (from %d to %d)", minPort, maxPort),
		Validate:  validatePort,
	}

	portString, err := prompt.Run()
	if (err != nil) {
		return 0, err
	}
	port, err = strconv.Atoi(portString)
	if (err != nil) {
		return 0, err
	}
	return port, nil
}

func validateString(name string, value string, minLength *int, maxLength *int) error {
	if (minLength != nil) {
		if len(value) == 0 {
			return fmt.Errorf("%s cannot be empty", name)
		} else if (len(value) < *minLength) {
			return fmt.Errorf("%s must be greater than %d", name, *minLength)
		}
	} else if (maxLength != nil && len(value) > *maxLength) {
		return fmt.Errorf("%s must be less than %d", name, *maxLength)
	}
	return nil
}

func getStringProperty(name string, minLength int, maxLength int) (string, error) {
	var property = viper.GetString(name)
	var err = validateString(name, property, &minLength, &maxLength)

	if  (err == nil) {
		return property, nil
	}
	valueDetails := fmt.Sprintf(": %s", property)
	if (property == "") {
		valueDetails = ""
	}
	showError(fmt.Sprintf("invalid config for %s%s", name, valueDetails), err)

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Enter %s", name),
		Validate:  func (value string) error { return validateString(name, value, &minLength, &maxLength) },
	}

	value, err := prompt.Run()
	if (err != nil) {
		return "", err
	}
	return value, nil
}

func authorize(clientID string, clientSecret string, port string, logger *zerolog.Logger) (*gaseed.GASeed, error) {
	server := gauth.New(clientID, clientSecret, port, logger)
	seed, err := server.Authorize(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	return seed, nil
}

func saveAccount(accountName string, seed *gaseed.GASeed) error {
	keyring := typedkeyring.New[gaseed.GASeed](serviceName)
	err := keyring.Save(accountName, seed)
	if err != nil {
		return fmt.Errorf("failed to save token %s to keyring: %w", accountName, err)
	}

	account, err := gaccount.New(serviceName, accountName)
	if err != nil {
		keyring.Delete(accountName)
		return fmt.Errorf("failed to add account '%s' to config: %w", accountName, err)
	}

	err = updateConfigFile(account)
	if err != nil {
		keyring.Delete(accountName)
		return fmt.Errorf("failed to add account '%s' to config: %w", accountName, err)
	}

	fmt.Printf("account '%s' was added to list of available accounts\n", accountName)
	return nil
}

func updateConfigFile(account *gaccount.GAccount) error {
	accounts := getAccountsFromConfig()
	accounts = append(accounts, *account)

	viper.Set(accountsConfigKey, accounts)
	return viper.WriteConfig()
}

func init() {
	addCmd.AddCommand(addAccountCmd)
	
	addAccountCmd.Flags().IntP("port", "p", defaultPort, "port number for gAuth code response")
	addAccountCmd.Flags().StringVar(&credFile, "credentials", "", "path to credentials file")

	viper.SetDefault("port", defaultPort)
	viper.BindPFlag("port", addAccountCmd.Flags().Lookup("port"))
}
