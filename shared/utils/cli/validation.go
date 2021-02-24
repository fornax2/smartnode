package cli

import (
    "fmt"
    "math/big"
    "regexp"
    "strconv"
    "strings"

    "github.com/ethereum/go-ethereum/common"
    "github.com/tyler-smith/go-bip39"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services/passwords"
)


// Config
const (
    MinDAOMemberIDLength = 3
    MinDAOMemberEmailLength = 6
)


//
// General types
//


// Validate command argument count
func ValidateArgCount(c *cli.Context, count int) error {
    if len(c.Args()) != count {
        return fmt.Errorf("Incorrect argument count; usage: %s", c.Command.UsageText)
    }
    return nil
}


// Validate a boolean value
func ValidateBool(name, value string) (bool, error) {
    val := strings.ToLower(value)
    if !(val == "true" || val == "yes" || val == "false" || val == "no") {
        return false, fmt.Errorf("Invalid %s '%s' - valid values are 'true', 'yes', 'false' and 'no'", name, value)
    }
    if val == "true" || val == "yes" {
        return true, nil
    } else {
        return false, nil
    }
}


// Validate an unsigned integer value
func ValidateUint(name, value string) (uint64, error) {
    val, err := strconv.ParseUint(value, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return val, nil
}


// Validate an address
func ValidateAddress(name, value string) (common.Address, error) {
    if !common.IsHexAddress(value) {
        return common.Address{}, fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return common.HexToAddress(value), nil
}


// Validate a wei amount
func ValidateWeiAmount(name, value string) (*big.Int, error) {
    val := new(big.Int)
    if _, ok := val.SetString(value, 10); !ok {
        return nil, fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return val, nil
}


// Validate an ether amount
func ValidateEthAmount(name, value string) (float64, error) {
    val, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return 0, fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return val, nil
}


// Validate a fraction
func ValidateFraction(name, value string) (float64, error) {
    val, err := strconv.ParseFloat(value, 64)
    if err != nil || val < 0 || val > 1 {
        return 0, fmt.Errorf("Invalid %s '%s' - must be a number between 0 and 1", name, value)
    }
    return val, nil
}


// Validate a percentage
func ValidatePercentage(name, value string) (float64, error) {
    val, err := strconv.ParseFloat(value, 64)
    if err != nil || val < 0 || val > 100 {
        return 0, fmt.Errorf("Invalid %s '%s' - must be a number between 0 and 100", name, value)
    }
    return val, nil
}


// Validate a token type
func ValidateTokenType(name, value string) (string, error) {
    val := strings.ToLower(value)
    if !(val == "eth" || val == "rpl" || val == "fsrpl" || val == "neth" || val == "reth") {
        return "", fmt.Errorf("Invalid %s '%s' - valid types are 'ETH', 'RPL', 'fsRPL', 'nETH' and 'rETH'", name, value)
    }
    return val, nil
}


//
// Command specific types
//


// Validate a positive wei amount
func ValidatePositiveWeiAmount(name, value string) (*big.Int, error) {
    val, err := ValidateWeiAmount(name, value)
    if err != nil {
        return nil, err
    }
    if val.Cmp(big.NewInt(0)) < 1 {
        return nil, fmt.Errorf("Invalid %s '%s' - must be greater than 0", name, value)
    }
    return val, nil
}


// Validate a deposit amount in wei
func ValidateDepositWeiAmount(name, value string) (*big.Int, error) {
    val, err := ValidateWeiAmount(name, value)
    if err != nil {
        return nil, err
    }
    if ether := strings.Repeat("0", 18); !(val.String() == "0" || val.String() == "16"+ether || val.String() == "32"+ether) {
        return nil, fmt.Errorf("Invalid %s '%s' - valid values are 0, 16 and 32 ether", name, value)
    }
    return val, nil
}


// Validate a positive ether amount
func ValidatePositiveEthAmount(name, value string) (float64, error) {
    val, err := ValidateEthAmount(name, value)
    if err != nil {
        return 0, err
    }
    if val <= 0 {
        return 0, fmt.Errorf("Invalid %s '%s' - must be greater than 0", name, value)
    }
    return val, nil
}


// Validate a deposit amount in ether
func ValidateDepositEthAmount(name, value string) (float64, error) {
    val, err := ValidateEthAmount(name, value)
    if err != nil {
        return 0, err
    }
    if !(val == 0 || val == 16 || val == 32) {
        return 0, fmt.Errorf("Invalid %s '%s' - valid values are 0, 16 and 32 ether", name, value)
    }
    return val, nil
}


// Validate a burnable token type
func ValidateBurnableTokenType(name, value string) (string, error) {
    val := strings.ToLower(value)
    if !(val == "neth" || val == "reth") {
        return "", fmt.Errorf("Invalid %s '%s' - valid types are 'nETH' and 'rETH'", name, value)
    }
    return val, nil
}


// Validate a node password
func ValidateNodePassword(name, value string) (string, error) {
    if len(value) < passwords.MinPasswordLength {
        return "", fmt.Errorf("Invalid %s '%s' - must be at least %d characters long", name, value, passwords.MinPasswordLength)
    }
    return value, nil
}


// Validate a wallet mnemonic phrase
func ValidateWalletMnemonic(name, value string) (string, error) {
    if !bip39.IsMnemonicValid(value) {
        return "", fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return value, nil
}


// Validate a timezone location
func ValidateTimezoneLocation(name, value string) (string, error) {
    if !regexp.MustCompile("^\\w{2,}\\/\\w{2,}$").MatchString(value) {
        return "", fmt.Errorf("Invalid %s '%s' - must be in the format 'Country/City'", name, value)
    }
    return value, nil
}


// Validate a DAO member ID
func ValidateDAOMemberID(name, value string) (string, error) {
    val := strings.TrimSpace(value)
    if len(val) < MinDAOMemberIDLength {
        return "", fmt.Errorf("Invalid %s '%s' - must be at least %d characters long", name, val, MinDAOMemberIDLength)
    }
    return val, nil
}


// Validate a DAO member email
func ValidateDAOMemberEmail(name, value string) (string, error) {
    if len(value) < MinDAOMemberEmailLength {
        return "", fmt.Errorf("Invalid %s '%s' - must be at least %d characters long", name, value, MinDAOMemberEmailLength)
    }
    if !regexp.MustCompile("^\\S+@\\S+(\\.\\S+)+$").MatchString(value) {
        return "", fmt.Errorf("Invalid %s '%s' - must be a valid email address", name, value)
    }
    return value, nil
}

