# Getting object permissions

> The sample shows how to get permissions for OData collections

Please be aware that `HasUniqueRoleAssignments` is a heavy property that creates workload on a SharePoint server.

Try not abusing it by potential requests to large lists getting a bunch of items.

```golang
// Getting all lists with role assigments details
res, err := sp.Web().Lists().
  Select(`
    Title,
    RoleAssignments/Member/*,
    RoleAssignments/RoleDefinitionBindings/*
  `).
  Expand(`
    RoleAssignments/Member,
    RoleAssignments/RoleDefinitionBindings
  `).
  Get()
```

`RoleAssigments` if any applied contains an array of objects.

```golang
var lists []*struct {
  Title           string
  RoleAssignments []*api.RoleAssigment
}

// .Normalized() method aligns responses between different OData modes
if err := json.Unmarshal(res.Normalized(), &lists); err != nil {
  log.Fatalf("unable to parse the response: %v", err)
}
```

Assigments is a Member and binded RoleDefinitions:

```golang
// RoleAssigment role asigments model
type RoleAssigment struct {
	Member *struct {
		LoginName     string
		PrincipalType int
	}
	RoleDefinitionBindings []*RoleDefInfo
}
```