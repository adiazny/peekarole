package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	filter string

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List ClusterRoles in the cluster",
		Long:  `List all ClusterRoles in the cluster with optional filtering`,
		RunE:  runList,
	}
)

func init() {
	listCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter ClusterRoles by name (supports partial matches)")
}

func runList(cmd *cobra.Command, args []string) error {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %v", err)
	}

	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing ClusterRoles: %v", err)
	}

	// tabwriter.NewWriter() method arguments:
	// - io.Writer: the writer to write to
	// - minwidth: the minimum width of each column
	// - tabwidth: the width of tab characters
	// - padding: the number of spaces between columns
	// - padchar: the character to use for padding
	// - flags: a bit field of options (none are currently defined)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "NAME\tRULES\tAGE\n")

	for _, cr := range clusterRoles.Items {
		if filter != "" && !strings.Contains(strings.ToLower(cr.Name), strings.ToLower(filter)) {
			continue
		}

		age := time.Since(cr.CreationTimestamp.Time).Round(time.Second)

		fmt.Fprintf(w, "%s\t%d\t%v\n", cr.Name, len(cr.Rules), age)
	}

	return w.Flush()
}
