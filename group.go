package main

import (
	"strings"
	"fmt"
)

type PackageGroupFunc func(category string) (packageNames []string)

func GroupEnumsAndFunctions(enums EnumCategories, functions FunctionCategories, sortFunc PackageGroupFunc) (packages Packages) {
	packages = make(Packages)
	for category, catEnums := range enums {
		if packageNames := sortFunc(category); packageNames != nil {
			for _, packageName := range packageNames {
				if p, ok := packages[packageName]; ok {
					p.Enums[category] = catEnums
				} else {
					p := &Package{make(EnumCategories), make(FunctionCategories)}
					p.Enums[category] = catEnums
					packages[packageName] = p
				}
			}
		}
	}
	for category, catFunctions := range functions {
		if packageNames := sortFunc(category); packageNames != nil {
			for _, packageName := range packageNames {
				if p, ok := packages[packageName]; ok {
					p.Functions[category] = catFunctions
				} else {
					p := &Package{make(EnumCategories), make(FunctionCategories)}
					p.Functions[category] = catFunctions
					packages[packageName] = p
				}
			}
		}
	}
	return
}

// Default grouping function:
func GroupPackagesByVendorFunc(category string, supportedVersions []Version) (packageNames []string) {
	pc, err := ParseCategoryString(category)
	if err != nil {
		return nil
	}
	packages := make([]string, 0, 8)
	switch pc.CategoryType {
		case CategoryExtension:
			packages = append(packages, strings.ToLower(pc.Vendor))
		case CategoryVersion:
			for _, ver := range supportedVersions {
				if pc.Version.Compare(ver) <= 0 {
					packages = append(packages, 
						fmt.Sprintf("gl%d%d", ver.Major, ver.Minor),
						fmt.Sprintf("gl%d%dd", ver.Major, ver.Minor))
				}
			}
		case CategoryDeprecatedVersion:
			for _, ver := range supportedVersions {
				if pc.Version.Compare(ver) <= 0 {
					packages = append(packages,
						fmt.Sprintf("gl%d%dd", ver.Major, ver.Minor))
				}
			}
	}
	return packages
}
