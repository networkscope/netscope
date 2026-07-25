package assets

import (
	"net"
	"strings"

	"github.com/networkscope/netscope/pkg/models"
)

func DiscoverTarget(target string) ([]*models.Asset, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, nil
	}

	var out []*models.Asset

	if ip := net.ParseIP(target); ip != nil {
		if ip.To4() != nil {
			a, err := models.NewAsset(target, models.AssetTypeIP, "input")
			if err != nil {
				return nil, err
			}
			out = append(out, a)
		} else {
			a, err := models.NewAsset(target, models.AssetTypeIP, "input")
			if err != nil {
				return nil, err
			}
			out = append(out, a)
		}
		return out, nil
	}

	if strings.Contains(target, ".") || strings.Contains(target, ":") {
		a, err := models.NewAsset(target, models.AssetTypeDomain, "input")
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	} else {
		a, err := models.NewAsset(target, models.AssetTypeHostname, "input")
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	return out, nil
}
