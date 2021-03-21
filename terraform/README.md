# terraform

IaC resources

1. requirements
    - ansible
    - terraform


# TL;DR;
```bash
terraform apply
```

# Posthumous setup
```bash
ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -u devbytom -i 'INSTANCE_IP,' --private-key ~/.ssh/id_rsa playbook.yml -e ansible_python_interpreter=/usr/bin/python3
```
