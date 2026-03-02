package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Simplified IR types — only the fields this generator needs.
type application struct {
	Name         string          `json:"name"`
	Config       *buildCfg       `json:"config,omitempty"`
	Architecture *architecture   `json:"architecture,omitempty"`
	Environments []*environment  `json:"environments,omitempty"`
	Data         []*model        `json:"data,omitempty"`
}

type buildCfg struct {
	Frontend string `json:"frontend"`
	Backend  string `json:"backend"`
	Database string `json:"database"`
	Deploy   string `json:"deploy"`
}

type architecture struct {
	Style    string    `json:"style"`
	Services []*svc    `json:"services,omitempty"`
}

type svc struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type environment struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars,omitempty"`
}

type model struct {
	Name string `json:"name"`
}

func runGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	irFile := fs.String("ir", "", "Path to IR JSON file")
	outputDir := fs.String("output", "", "Output directory")
	settingsStr := fs.String("settings", "", "Plugin settings JSON")
	fs.Parse(args)

	if *irFile == "" || *outputDir == "" {
		return fmt.Errorf("--ir and --output are required")
	}

	data, err := os.ReadFile(*irFile)
	if err != nil {
		return fmt.Errorf("reading IR: %w", err)
	}

	var app application
	if err := json.Unmarshal(data, &app); err != nil {
		return fmt.Errorf("parsing IR: %w", err)
	}

	var settings map[string]string
	if *settingsStr != "" {
		json.Unmarshal([]byte(*settingsStr), &settings)
	}

	return generate(&app, *outputDir, settings)
}

func generate(app *application, outputDir string, settings map[string]string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	namespace := settings["namespace"]
	if namespace == "" {
		namespace = strings.ToLower(app.Name)
	}

	files := map[string]string{
		filepath.Join(outputDir, "namespace.yaml"):  generateNamespace(namespace),
		filepath.Join(outputDir, "deployment.yaml"): generateDeployment(app, namespace),
		filepath.Join(outputDir, "service.yaml"):    generateService(app, namespace),
		filepath.Join(outputDir, "ingress.yaml"):    generateIngress(app, namespace),
	}

	// Generate configmaps for each environment.
	for _, env := range app.Environments {
		filename := fmt.Sprintf("configmap-%s.yaml", strings.ToLower(env.Name))
		files[filepath.Join(outputDir, filename)] = generateConfigMap(app, namespace, env)
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func generateNamespace(namespace string) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Namespace
metadata:
  name: %s
  labels:
    app.kubernetes.io/managed-by: human-compiler
`, namespace)
}

func generateDeployment(app *application, namespace string) string {
	name := strings.ToLower(app.Name)
	replicas := 2
	port := 3000

	// Determine container image and port from config.
	image := fmt.Sprintf("%s:latest", name)
	if app.Config != nil {
		if strings.Contains(strings.ToLower(app.Config.Backend), "go") {
			port = 8080
		} else if strings.Contains(strings.ToLower(app.Config.Backend), "python") {
			port = 8000
		}
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
spec:
  replicas: %d
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: %s
        image: %s
        ports:
        - containerPort: %d
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: %d
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: %d
          initialDelaySeconds: 5
          periodSeconds: 10
        envFrom:
        - configMapRef:
            name: %s-config
`, name, namespace, name, replicas, name, name, name, image, port, port, port, name))

	// Add database connection if configured.
	if app.Config != nil && app.Config.Database != "" {
		b.WriteString(fmt.Sprintf(`        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: %s-secrets
              key: database-url
`, name))
	}

	return b.String()
}

func generateService(app *application, namespace string) string {
	name := strings.ToLower(app.Name)
	port := 3000
	if app.Config != nil {
		if strings.Contains(strings.ToLower(app.Config.Backend), "go") {
			port = 8080
		} else if strings.Contains(strings.ToLower(app.Config.Backend), "python") {
			port = 8000
		}
	}

	return fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: %d
    protocol: TCP
  selector:
    app: %s
`, name, namespace, name, port, name)
}

func generateIngress(app *application, namespace string) string {
	name := strings.ToLower(app.Name)
	host := fmt.Sprintf("%s.example.com", name)

	return fmt.Sprintf(`apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s
  namespace: %s
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: %s
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: %s
            port:
              number: 80
  tls:
  - hosts:
    - %s
    secretName: %s-tls
`, name, namespace, host, name, host, name)
}

func generateConfigMap(app *application, namespace string, env *environment) string {
	name := strings.ToLower(app.Name)
	envName := strings.ToLower(env.Name)

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`apiVersion: v1
kind: ConfigMap
metadata:
  name: %s-config
  namespace: %s
  labels:
    app: %s
    environment: %s
data:
`, name, namespace, name, envName))

	for k, v := range env.Vars {
		b.WriteString(fmt.Sprintf("  %s: %q\n", k, v))
	}

	return b.String()
}
