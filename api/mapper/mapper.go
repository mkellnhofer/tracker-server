package mapper

import (
	aModel "kellnhofer.com/tracker/api/model"
	lModel "kellnhofer.com/tracker/model"
)

func ToApiLocs(iLocs []*lModel.Location) []*aModel.Location {
	oLocs := []*aModel.Location{}
	for _, iLoc := range iLocs {
		oLocs = append(oLocs, ToApiLoc(iLoc))
	}
	return oLocs
}

func ToApiLoc(iLoc *lModel.Location) *aModel.Location {
	return &aModel.Location{iLoc.Id, iLoc.ChangeTime, iLoc.Name, iLoc.Time, iLoc.Lat, iLoc.Lng,
		iLoc.Description, ToApiPers(iLoc.Persons)}
}

func ToApiPers(iPers []*lModel.Person) []*aModel.Person {
	oPers := []*aModel.Person{}
	for _, iPer := range iPers {
		oPers = append(oPers, ToApiPer(iPer))
	}
	return oPers
}

func ToApiPer(iPer *lModel.Person) *aModel.Person {
	return &aModel.Person{iPer.FirstName, iPer.LastName}
}

func ToLogicLoc(iLoc *aModel.Location) *lModel.Location {
	return &lModel.Location{iLoc.Id, 0, iLoc.Name, iLoc.Time, iLoc.Lat, iLoc.Lng, iLoc.Description,
		ToLogicPers(iLoc.Persons)}
}

func ToLogicPers(iPers []*aModel.Person) []*lModel.Person {
	if iPers == nil {
		return nil
	}

	var oPers []*lModel.Person
	for _, iPer := range iPers {
		oPers = append(oPers, ToLogicPer(iPer))
	}
	return oPers
}

func ToLogicPer(iPer *aModel.Person) *lModel.Person {
	return &lModel.Person{0, iPer.FirstName, iPer.LastName}
}
