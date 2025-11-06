package provider

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/larstobi/go-multipass/multipass"
)

func QueryInstance(state Instance) (*Instance, error) {
	out, err := multipass.Get(&multipass.GetReq{Name: state.Name.Value})
	if err != nil {
		return nil, errors.New("Could not get local value. " + err.Error())
	}

	// Check CPUS
	current_cpus := state.CPUS
	// If CPUS is not specified, then ignore changes
	if !state.CPUS.Null {
		cpus := new(big.Float)
		cpus.SetString(out.CPUS)
		if cpus != state.CPUS.Value {
			current_cpus = types.Number{Value: cpus}
		}
	}

	// Check Memory
	current_memory := state.Memory
	// If Memory is not specified, then ignore changes
	if !state.Memory.Null {
		if equal, err := CompareDataSizes(out.Memory, state.Memory.Value); err != nil {
			return nil, errors.New("Error comparing memory size: " + err.Error())
		} else {
			if !*equal {
				current_memory = types.String{Value: RemoveZeroDecimalAndSpaces(out.Memory)}
			}
		}
	}

	// Check Disk
	current_disk := state.Disk
	// If Disk is not specified, then ignore changes
	if !state.Disk.Null {
		if equal, err := CompareDataSizes(out.Disk, state.Disk.Value); err != nil {
			return nil, errors.New("Error comparing disk size: " + err.Error())
		} else {
			if !*equal {
				current_disk = types.String{Value: RemoveZeroDecimalAndSpaces(out.Disk)}
			}
		}
	}

	// Fetch instance info to get IP address
	instanceInfo, infoErr := multipass.Info(&multipass.InfoRequest{Name: state.Name.Value})
	ipv4 := state.IPv4
	if infoErr == nil && instanceInfo != nil {
		ipv4 = types.String{Value: instanceInfo.IP}
	}

	// Generate resource state struct
	var result = Instance{
		Name:          state.Name,
		Image:         state.Image,
		CPUS:          current_cpus,
		Memory:        current_memory,
		Disk:          current_disk,
		CloudInitFile: state.CloudInitFile,
		Network:       state.Network,
		Bridged:       state.Bridged,
		IPv4:          ipv4,
	}

	return &result, nil
}

func ToMegabytes(input string) (string, error) {
	input = RemoveZeroDecimalAndSpaces(input)
	input = strings.ReplaceAll(input, "\"", "")

	output, err := bytefmt.ToMegabytes(input)
	if err != nil {
		return "", err

	}

	result := strconv.FormatUint(output, 10)
	return result, nil
}

// multipass get returns memory and disk values with a .0 decimal
func RemoveZeroDecimalAndSpaces(input string) string {
	input = strings.ReplaceAll(input, " ", "")
	sep := strings.LastIndexAny(input, "01234567890. ")
	var number, suffix string
	number = input[:sep+1]
	suffix = input[sep+1:]
	return strings.TrimSuffix(number, ".0") + suffix
}

func CompareDataSizes(value1 string, value2 string) (*bool, error) {
	result1, result1_err := ToMegabytes(value1)
	if result1_err != nil {
		return nil, result1_err
	}

	result2, result2_err := ToMegabytes(value2)
	if result2_err != nil {
		return nil, result2_err
	}

	equal := result1 == result2
	return &equal, nil
}
