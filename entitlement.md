# Entitlement

## Model
The draft model to save profile/entitlement data

```golang
type Entitlement struct {
	Slug string `json:"slug" firestore:"slug" bigquery:"slug"`
}

type EntitlementProfile struct {
	Slug         string        `json:"slug" firestore:"slug" bigquery:"slug"`
	Entitlements []Entitlement `json:"entitlements" firestore:"entitlements" bigquery:"entitlements"`
	Data         string        `json:"-" firestore:"-" bigquery:"data"`
}
```

## Profiles
Will be saved to both DBs (Firestore and Bigquery)

### Structure
collection: `entitlementProfile`

entries:
```json
{
    "admin:": {
        "slug": "admin",
        "entitlements": [{
            "slug": "broker.emit"
        }, {
            "slug": "broker.lead"
        }],
        "data": "{\"slug\": \"admin\", \"entitlements\": [{\"slug\": \"broker.emit\"}, {\"slug\": \"broker.lead\"}]}"
    }
}
```

## Changes to NetworkNode
add new field `entitlementProfile` that points to the profile name and possible overwrites

```golang
// TODO: add column to BigQuery Schema
type NetworkNode struct {
    // ...
    EntitlementProfile  EntitlementProfile    `json:"entitlementProfile,omitempty" firestore:"entitlementProfile,omitempty" bigquery:"entitlementProfile"`
}
```

Example:
```jsonc
// No reference - gets entitlement by authtoken.role (retrocompatibility with current nodes)
{
    // ...
    "entitlementProfile": {
        "slug": "",
        "entitlements": [],
        // data is not present since it is saved only in bigquery for safekeeping
    }
}
// With profile
{
    // ...
    "entitlementProfile": {
        "slug": "my-custom-or-not-profile",
        "entitlements": [],
        // data is not present since it is saved only in bigquery for safekeeping
    }
}
// With overrides
{
    // ...
    "entitlementProfile": {
        "slug": "the-profile-does-not-matter",
        "entitlements": [{
            "slug": "one"
        }, {
            "slug": "two"
        }],
        // data is not present since it is saved only in bigquery for safekeeping
    }
}
```

## Changes to Get NetworkNode(s)
The `EntitlementProfile` field must be filled with one of the three possible scenarios:

1. Nothing set -> get data from `entitlementProfile` collection by `authToken.role`
2. Profile set -> get data from `entitlementProfile` collection by `entitlementProfile.slug`
3. Override -> use the data already presented on the node

```golang
nn := models.NetworkNode{}
authToken := lib.AuthToken{}
if len(nn.EntitlementProfile.Entitlements) > 0 {
    return
}
if nn.EntitlementProfile.Slug != "" {
    nn.EntitlementProfile.Entitlements := GetFromFirestore(nn.EntitlementProfile.Slug)
    return
}
nn.EntitlementProfile.Entitlements := GetFromFirestore(authToken.Role)
```

## Changes to User
there are ate least two possibilities:
- admins keep being users:
  - the same strategy that was added to the networkNodes is implemented for users, so the frontend gets the values
  - the frontend invokes a new endpoint that retrieves the entitlement by slug (should use the authtoken.role)
- admins become networkNodes:
  - needs a discovery fo all possible impacts other than recreating all admin users

The one with the lowest effort seems the second where in the case of a non networkNode the frontend makes another call to get the entitlement data. This should also work for customers on the client app (non back office)