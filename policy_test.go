package policy

import "testing"

func TestRules(t *testing.T) {

	context := Context{
		Roles: []string{"guest", "member"},
		Token: map[string]string{
			"user_id":    "u-1",
			"project_id": "p-2",
		},
		Request: map[string]string{
			"target.user_id": "u-1",
			"user_id":        "u-2",
		},
	}

	testCases := []struct {
		rule   string
		result bool
	}{
		{"", true},
		{"@", true},
		{"!", false},
		{"role:member", true},
		{"not role:member", false},
		{"role:admin", false},
		{"role:admin or role:guest", true},
		{"role:admin and role:guest", false},
		{"user_id:u-1", true},
		{"user_id:u-2", false},
		{"'u-2':%(user_id)s", true},
		{"domain_id:%(does_not_exit)s", false},
		{"not (@ or @)", false},
		{"not @ or @", true},
		{"@ and (! or (not !))", true},
	}

	for _, c := range testCases {
		//fmt.Println("Testing rule ", c.rule)
		rule, err := parseRule(c.rule)
		if err != nil {
			t.Errorf("Failed to parse rule %q: %s", c.rule, err)
			continue
		}
		if result := rule(context); result != c.result {
			t.Errorf("Rule %q returned %v, expected %v", c.rule, result, c.result)
		}
	}

}

func TestPolicy(t *testing.T) {
	testPolicy := map[string]string{
		"admin_required":   "role:admin",
		"service_role":     "role:service",
		"service_or_admin": "rule:admin_required or rule:service_role",
	}
	serviceContext := Context{
		Roles: []string{"service"},
	}

	enforcer, err := NewPolicy(testPolicy)
	if err != nil {
		t.Fatal("Failed to parse policy ", err)
	}
	if !enforcer.Enforce("service_or_admin", serviceContext) {
		t.Error("service_or_admin check should have returned true")
	}

}
