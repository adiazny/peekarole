package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var compareCmd = &cobra.Command{
	Use:   "compare [clusterrole1] [clusterrole2]",
	Short: "Compare two or more ClusterRoles",
	Long:  `Compare two or more ClusterRoles to see their permission differences`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runCompare,
}

func runCompare(cmd *cobra.Command, args []string) error {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %v", err)
	}

	clusterRoles := make(map[string]*rbacv1.ClusterRole)
	for _, name := range args {
		cr, err := clientset.RbacV1().ClusterRoles().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting ClusterRole %s: %v", name, err)
		}
		clusterRoles[name] = cr
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "COMPARISON RESULTS:")
	fmt.Fprintln(w, "==================")

	// Compare number of rules
	fmt.Fprintln(w, "\nRule Count:")
	for name, cr := range clusterRoles {
		fmt.Fprintf(w, "%s:\t%d rules\n", name, len(cr.Rules))
	}

	// Show common permissions first
	fmt.Fprintln(w, "\nCommon Permissions:")
	commonPerms := getCommonPermissions(clusterRoles)
	if len(commonPerms) > 0 {
		for _, perm := range commonPerms {
			fmt.Fprintf(w, "  - %s\n", formatPermission(perm))
		}
	} else {
		fmt.Fprintln(w, "  No common permissions found")
	}

	// Compare unique permissions
	fmt.Fprintln(w, "\nUnique Permissions:")
	for name, cr := range clusterRoles {
		uniquePerms := getUniquePermissions(cr, clusterRoles)
		if len(uniquePerms) > 0 {
			fmt.Fprintf(w, "\n%s unique permissions:\n", name)
			for _, perm := range uniquePerms {
				fmt.Fprintf(w, "  - %s\n", formatPermission(perm))
			}
		} else {
			fmt.Fprintf(w, "\n%s has no unique permissions\n", name)
		}
	}

	return w.Flush()
}

type Permission struct {
	Verbs     []string
	APIGroups []string
	Resources []string
}

func getUniquePermissions(cr *rbacv1.ClusterRole, allRoles map[string]*rbacv1.ClusterRole) []Permission {
	var unique []Permission

	for _, rule := range cr.Rules {
		perm := Permission{
			Verbs:     rule.Verbs,
			APIGroups: rule.APIGroups,
			Resources: rule.Resources,
		}

		isUnique := true
		for name, otherCR := range allRoles {
			if name == cr.Name {
				continue
			}

			if hasPermission(otherCR, perm) {
				isUnique = false
				break
			}
		}

		if isUnique {
			unique = append(unique, perm)
		}
	}

	return unique
}

func getCommonPermissions(roles map[string]*rbacv1.ClusterRole) []Permission {
	if len(roles) < 2 {
		return nil
	}

	// Get all permissions from the first role as initial common set
	var firstRole *rbacv1.ClusterRole
	for _, role := range roles {
		firstRole = role
		break
	}

	common := make([]Permission, 0)
	for _, rule := range firstRole.Rules {
		perm := Permission{
			Verbs:     rule.Verbs,
			APIGroups: rule.APIGroups,
			Resources: rule.Resources,
		}

		// Check if this permission exists in all other roles
		isCommon := true
		for name, otherRole := range roles {
			if name == firstRole.Name {
				continue
			}
			if !hasPermission(otherRole, perm) {
				isCommon = false
				break
			}
		}

		if isCommon {
			common = append(common, perm)
		}
	}

	return common
}

func hasPermission(cr *rbacv1.ClusterRole, perm Permission) bool {
	for _, rule := range cr.Rules {
		if containsAll(rule.Verbs, perm.Verbs) &&
			containsAll(rule.APIGroups, perm.APIGroups) &&
			containsAll(rule.Resources, perm.Resources) {
			return true
		}
	}
	return false
}

func containsAll(haystack, needles []string) bool {
	for _, needle := range needles {
		found := false
		for _, hay := range haystack {
			if hay == needle || hay == "*" {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func formatPermission(p Permission) string {
	sort.Strings(p.Verbs)
	sort.Strings(p.APIGroups)
	sort.Strings(p.Resources)

	apiGroups := "*"
	if len(p.APIGroups) > 0 && p.APIGroups[0] != "" {
		apiGroups = strings.Join(p.APIGroups, ",")
	}

	return fmt.Sprintf("[%s] %s: %s",
		apiGroups,
		strings.Join(p.Resources, ","),
		strings.Join(p.Verbs, ","))
}
