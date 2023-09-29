package adapters

import (
	"net/url"

	"github.com/mrparano1d/noxite/ent"
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
