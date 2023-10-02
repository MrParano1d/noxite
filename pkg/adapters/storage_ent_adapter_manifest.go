package adapters

import (
	"fmt"
	"net/url"

	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

func bugsFromFieldBugs(bgs *fields.Bugs) *bugs {
	if bgs == nil {
		return nil
	}
	return &bugs{
		URL:   bgs.URL.String(),
		Email: bgs.Email.String(),
	}
}

func authorFromMixedAuthor(pkgAuthor *fields.MixedAuthor) *author {
	if pkgAuthor == nil {
		return nil
	}
	// TODO: implement with EntityID Author
	return nil
}

func contributorsFromMixedAuthors(pkgContributors fields.MixedAuthors) []contributor {
	if pkgContributors == nil {
		return nil
	}
	// TODO: implement with EntityID Author
	return nil
}

func stringMapFromRequiredStringMap(m map[fields.RequiredString]fields.RequiredString) map[string]string {
	if m == nil {
		return nil
	}
	sm := make(map[string]string)
	for k, v := range m {
		sm[k.String()] = v.String()
	}
	return sm
}

func directoriesFromFieldDirectories(dirs *fields.Directories) *directories {
	if dirs == nil {
		return nil
	}

	return &directories{
		Bin: dirs.Bin.String(),
		Man: dirs.Man.String(),
	}
}

func repositoryFromFrieldRepository(repo *fields.Repository) *repository {
	if repo == nil {
		return nil
	}

	typ := new(string)

	if repo.Type != nil {
		*typ = repo.Type.String()
	}

	return &repository{
		Type: typ,
		URL:  repo.URL.String(),
	}
}

func peerDependenciesMetaFromFieldPeerDependenciesMeta(meta map[fields.RequiredString]map[fields.RequiredString]any) map[string]map[string]any {
	if meta == nil {
		return nil
	}

	pdm := make(map[string]map[string]any)
	for k, v := range meta {
		pdm[k.String()] = make(map[string]any)
		for k2, v2 := range v {
			pdm[k.String()][k2.String()] = v2
		}
	}
	return pdm
}

func mapAnyFromMapRequiredStringAny(m map[fields.RequiredString]any) map[string]any {
	if m == nil {
		return nil
	}

	mm := make(map[string]any)
	for k, v := range m {
		mm[k.String()] = v
	}
	return mm
}

func fundingFromMixedFunding(funding []fields.UrlType) any {
	if funding == nil {
		return nil
	}

	var f []any
	for _, v := range funding {
		if v.URL == nil {
			continue
		}
		if v.Type == nil {
			f = append(f, v.URL.String())
		} else {
			f = append(f, struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			}{
				URL:  v.URL.String(),
				Type: v.Type.String(),
			})
		}
	}

	return f
}

func manifestFromPackageVersion(packageName fields.RequiredString, ver *entities.PackageVersion) manifest {

	var description string
	if ver.Description != nil {
		description = *ver.Description
	}
	var license string
	if ver.License != nil {
		license = *ver.License
	}
	var main string
	if ver.Main != nil {
		main = ver.Main.String()
	}
	var browser string
	if ver.Browser != nil {
		browser = ver.Browser.String()
	}
	var readme string
	if ver.Readme != nil {
		readme = *ver.Readme
	}

	versions := make(map[string]revision)
	versions[ver.Version.String()] = revision{
		Name:                 packageName.String(),
		Version:              ver.Version.String(),
		Description:          description,
		Keywords:             fields.StringsFromRequiredStrings(ver.Keywords),
		Homepage:             ver.Homepage.String(),
		Bugs:                 bugsFromFieldBugs(ver.Bugs),
		License:              license,
		Author:               authorFromMixedAuthor(ver.Author),
		Contributors:         contributorsFromMixedAuthors(ver.Contributors),
		Funding:              fundingFromMixedFunding(ver.Funding),
		Files:                fields.StringsFromRequiredStrings(ver.Files),
		Main:                 main,
		Browser:              browser,
		Bin:                  stringMapFromRequiredStringMap(ver.Bin),
		Man:                  fields.StringsFromRequiredStrings(ver.Man),
		Directories:          directoriesFromFieldDirectories(ver.Directories),
		Repository:           repositoryFromFrieldRepository(ver.Repository),
		Scripts:              stringMapFromRequiredStringMap(ver.Scripts),
		Config:               stringMapFromRequiredStringMap(ver.Config),
		Dependencies:         stringMapFromRequiredStringMap(ver.Dependencies),
		DevDependencies:      stringMapFromRequiredStringMap(ver.DevDependencies),
		PeerDependencies:     stringMapFromRequiredStringMap(ver.PeerDependencies),
		PeerDependenciesMeta: peerDependenciesMetaFromFieldPeerDependenciesMeta(ver.PeerDependenciesMeta),
		BundledDependencies:  fields.StringsFromRequiredStrings(ver.BundledDependencies),
		OptionalDependencies: stringMapFromRequiredStringMap(ver.OptionalDependencies),
		Engines:              stringMapFromRequiredStringMap(ver.Engines),
		OS:                   fields.StringsFromRequiredStrings(ver.OS),
		CPU:                  fields.StringsFromRequiredStrings(ver.CPU),
		Private:              ver.Private,
		PublishConfig:        mapAnyFromMapRequiredStringAny(ver.PublishConfig),
		Workspaces:           fields.StringsFromRequiredStrings(ver.Workspaces),
		Dist: dist{
			Tarball:   "http://localhost:3000/" + url.QueryEscape(packageName.String()) + "/-/" + url.QueryEscape(packageName.String()) + "-" + ver.Version.String() + ".tgz",
			Integrity: ver.Integrity.String(),
			SHASUM:    ver.SHASUM.String(),
		},
	}

	attachments := make(map[string]attachment)
	attachments[ver.Version.String()] = attachment{
		ContentType: ver.ContentType.String(),
		Data:        ver.Data.String(),
		Length:      ver.Length,
	}
	m := manifest{
		Name:        packageName.String(),
		Description: description,
		Readme:      readme,
		Versions:    versions,
		Attachments: attachments,
		DistTags:    map[string]string{},
	}

	return m

}

func manifestFromEntVersion(packageName string, ver *ent.Version) manifest {
	versions := make(map[string]revision)
	versions[ver.Version] = revision{
		Name:                 packageName,
		Version:              ver.Version,
		Description:          ver.Description,
		Keywords:             ver.Keywords,
		Homepage:             ver.Homepage,
		Bugs:                 bugsFromFieldBugs(ver.Bugs),
		License:              ver.License,
		Author:               authorFromMixedAuthor(ver.Author),
		Contributors:         contributorsFromMixedAuthors(ver.Contributors),
		Funding:              fundingFromMixedFunding(ver.Funding),
		Files:                ver.Files,
		Main:                 ver.Main,
		Browser:              ver.Browser,
		Bin:                  stringMapFromRequiredStringMap(ver.Bin),
		Man:                  fields.StringsFromRequiredStrings(ver.Man),
		Directories:          directoriesFromFieldDirectories(ver.Directories),
		Repository:           repositoryFromFrieldRepository(ver.Repository),
		Scripts:              stringMapFromRequiredStringMap(ver.Scripts),
		Config:               stringMapFromRequiredStringMap(ver.Config),
		Dependencies:         stringMapFromRequiredStringMap(ver.Dependencies),
		DevDependencies:      stringMapFromRequiredStringMap(ver.DevDependencies),
		PeerDependencies:     stringMapFromRequiredStringMap(ver.PeerDependencies),
		PeerDependenciesMeta: peerDependenciesMetaFromFieldPeerDependenciesMeta(ver.PeerDependenciesMeta),
		BundledDependencies:  ver.BundledDependencies,
		OptionalDependencies: stringMapFromRequiredStringMap(ver.OptionalDependencies),
		Engines:              stringMapFromRequiredStringMap(ver.Engines),
		OS:                   ver.Os,
		CPU:                  ver.CPU,
		Private:              &ver.Private,
		PublishConfig:        mapAnyFromMapRequiredStringAny(ver.PublishConfig),
		Workspaces:           ver.Workspaces,
		Dist: dist{
			Tarball: "http://localhost:3000/" + url.QueryEscape(packageName) + "/-/" + url.QueryEscape(packageName) + "-" + ver.Version + ".tgz",

			Integrity: ver.Integrity,
			SHASUM:    ver.Shasum,
		},
	}

	attachments := make(map[string]attachment)
	attachments[ver.Version] = attachment{
		ContentType: ver.ContentType,
		Data:        ver.Data,
		Length:      ver.Length,
	}
	m := manifest{
		Name:        packageName,
		Description: ver.Description,
		Readme:      ver.Readme,
		Versions:    versions,
		Attachments: attachments,
		DistTags:    map[string]string{},
	}

	return m
}

func packageVersionFromEntVersion(packageName string, ver *ent.Version) (*entities.PackageVersion, error) {

	name, err := fields.RequiredStringFromString(packageName)
	if err != nil {
		return nil, &InvalidPackageVersionFieldErrror{Field: "name", Rearson: err.Error()}
	}

	version, err := fields.RequiredStringFromString(ver.Version)
	if err != nil {
		return nil, &InvalidPackageVersionFieldErrror{Field: "version", Rearson: err.Error()}
	}

	var description *string
	if ver.Description != "" {
		description = &ver.Description
	}

	keywords := make([]fields.RequiredString, len(ver.Keywords))
	for i, v := range ver.Keywords {
		keywords[i], err = fields.RequiredStringFromString(v)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "keywords", Rearson: err.Error()}
		}
	}

	var license *string
	if ver.License != "" {
		license = &ver.License
	}

	var main fields.RequiredString
	if ver.Main != "" {
		main, err = fields.RequiredStringFromString(ver.Main)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "main", Rearson: err.Error()}
		}
	}

	var browser fields.RequiredString
	if ver.Browser != "" {
		browser, err = fields.RequiredStringFromString(ver.Browser)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "browser", Rearson: err.Error()}
		}
	}

	var readme *string
	if ver.Readme != "" {
		readme = &ver.Readme
	}

	var homepage fields.Website
	if ver.Homepage != "" {
		homepage, err = fields.WebsiteFromString(ver.Homepage)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "homepage", Rearson: err.Error()}
		}
	}

	files := make([]fields.RequiredString, len(ver.Files))
	for i, v := range ver.Files {
		files[i], err = fields.RequiredStringFromString(v)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "files", Rearson: err.Error()}
		}
	}

	bundledDependencies := make([]fields.RequiredString, len(ver.BundledDependencies))
	for i, v := range ver.BundledDependencies {
		bundledDependencies[i], err = fields.RequiredStringFromString(v)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "bundledDependencies", Rearson: err.Error()}
		}
	}

	os := make([]fields.RequiredString, len(ver.Os))
	for i, v := range ver.Os {
		os[i], err = fields.RequiredStringFromString(v)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "os", Rearson: err.Error()}
		}
	}

	cpu := make([]fields.RequiredString, len(ver.CPU))
	for i, v := range ver.CPU {
		cpu[i], err = fields.RequiredStringFromString(v)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "cpu", Rearson: err.Error()}
		}
	}

	worspaces := make([]fields.RequiredString, len(ver.Workspaces))
	for i, v := range ver.Workspaces {
		worspaces[i], err = fields.RequiredStringFromString(v)
		if err != nil {
			return nil, &InvalidPackageVersionFieldErrror{Field: "workspaces", Rearson: err.Error()}
		}
	}

	integrity, err := fields.RequiredStringFromString(ver.Integrity)
	if err != nil {
		return nil, &InvalidPackageVersionFieldErrror{Field: "integrity", Rearson: err.Error()}
	}

	shasum, err := fields.RequiredStringFromString(ver.Shasum)
	if err != nil {
		return nil, &InvalidPackageVersionFieldErrror{Field: "shasum", Rearson: err.Error()}
	}

	contentType, err := fields.RequiredStringFromString(ver.ContentType)
	if err != nil {
		return nil, &InvalidPackageVersionFieldErrror{Field: "contentType", Rearson: err.Error()}
	}

	data, err := fields.RequiredStringFromString(ver.Data)
	if err != nil {
		return nil, &InvalidPackageVersionFieldErrror{Field: "data", Rearson: err.Error()}
	}

	return &entities.PackageVersion{
		Name:                 name,
		Version:              version,
		Description:          description,
		Keywords:             keywords,
		Homepage:             &homepage,
		Bugs:                 ver.Bugs,
		License:              license,
		Author:               ver.Author,
		Contributors:         ver.Contributors,
		Funding:              ver.Funding,
		Files:                files,
		Main:                 &main,
		Browser:              &browser,
		Bin:                  ver.Bin,
		Man:                  ver.Man,
		Directories:          ver.Directories,
		Repository:           ver.Repository,
		Scripts:              ver.Scripts,
		Config:               ver.Config,
		Dependencies:         ver.Dependencies,
		DevDependencies:      ver.DevDependencies,
		PeerDependencies:     ver.PeerDependencies,
		PeerDependenciesMeta: ver.PeerDependenciesMeta,
		BundledDependencies:  bundledDependencies,
		OptionalDependencies: ver.OptionalDependencies,
		Engines:              ver.Engines,
		OS:                   os,
		CPU:                  cpu,
		Private:              &ver.Private,
		PublishConfig:        ver.PublishConfig,
		Workspaces:           worspaces,

		Integrity:   integrity,
		SHASUM:      shasum,
		ContentType: contentType,
		Data:        data,
		Length:      ver.Length,
		Readme:      readme,
	}, nil

}

type InvalidPackageVersionFieldErrror struct {
	Field   string
	Rearson string
}

func (e *InvalidPackageVersionFieldErrror) Error() string {
	return fmt.Sprintf("invalid package version field %s: %s", e.Field, e.Rearson)
}
