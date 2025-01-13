package config

import (
    "os"
    "path/filepath"
)

var (
    CAFile = configFile("ca.pem")
    ServerCertFile = configFile("server.pem")
    ServerKeyFile = configFile("server-key.pem")
    RootClientCertFile = configFile("root-client.pem")
    RootClientKeyFile = configFile("root-client-key.pem")
    NobodyClientCertFile = configFile("nobody-client.pem")
    NobodyClientKeyFile = configFile("nobody-client-key.pem")
    ACLModelFile = configFile("model.conf")
    ACLPolicyFile = configFile("policy.csv")
)

func configFile(filename string) string {
    if dir := os.Getenv("CONFIG_DIR"); dir != "" {
        return filepath.Join(dir, filename)
    }
    // use $HOME/.proglog if $CONFIG_DIR isn't set
    dir, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }
    return filepath.Join(dir, ".proglog", filename)
}

