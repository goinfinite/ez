```
This project is under active development and is not ready for production use.
```

# Speedia FleetManager

Speedia FleetManager (SFM) is proprietary container management platform in a single file. It has a REST API, CLI and dashboard.

## Running

## Development

In this repository you'll find the REST API and CLI code plus the dashboard assets. The API and CLI uses Clean Architecture, DDD, TDD, CQRS, Object Calisthenics, etc. Understand how these concepts works before proceeding is advised.

To run this project during development you must install [Air](https://github.com/cosmtrek/air). Air is a tool that will watch for changes in the project and recompile it automatically.

### Unit Testing

Since SFM relies on the operational system being openSUSE MicroOS, the entire development and testing should be done in a VM. The VM can be created with the following steps:

1. Install VMWare Player;

2. Download the VMware `.vmx` and `.vmdk` files from "Base System + Container Runtime" column on MicroOS download page:
   https://en.opensuse.org/Portal:MicroOS/Downloads

3. Add the VM to the VMWare Player interface, change the resources to 2GB RAM // 2vCPU, fix the disk path and add a Network Adapter so the VM gets a connection.

4. Run the VM. The first boot will allow you to set up a root password. After the first reboot, login with the password you set.

5. On the terminal, change "security=1 selinux=1" to "selinux=0" with the command:

```
vim /etc/default/grub
```

Then apply the changes with the following command (no need to reboot yet):

```
transactional-update grub.cfg
```

6. Install a few system packages and Go runtime:

```
transactional-update pkg install git wget curl cyrus-sasl pam-devel gcc make tar procps
wget -nv https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.bashrc
echo 'alias sfm-swag="swag init -g src/presentation/api/api.go -o src/presentation/api/docs"' >> ~/.bashrc
```

7. Add your public SSH key to "/root/.ssh/authorized_keys" file and now you can reboot:

```
transactional-update reboot
```

8. After the reboot, you can install a few Go packages and clone the SFM repository:

```
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/cosmtrek/air@latest
git config --global user.name "yourgithubnick"
git config --global user.email yourgithubemail
git clone git@github.com:speedianet/sfm.git
```

OK, now you have the VM running and you should be able to use the [Visual Studio Remote SSH extension](https://code.visualstudio.com/docs/remote/ssh) to connect to it and manage the project.

Make sure to use the SSH key to connect to the VM and not the password. The IP address of the VM can be found with the `ip a` command on the VM terminal.

### Environment Variables

You must have an `.env` file in the root of the git directory **during development**. You can use the `.env.example` file as a template. Air will read the `.env` file and use it to run the project during development.

If you add a new env var that is required to run the apis, please add it to the `src/presentation/shared/checkEnvs.go` file.

When running in production, the `/speedia/.env` file is only used if the environment variables weren't set in the system. For instance, if you want to set the `ENV1` variable, you can do it in the `.env` file or in the command line:

```
ENV1=XXX /speedia/sfm
```

### Dev Utils

The `src/devUtils` folder is not a Clean Architecture layer, it's there to help you during development. You can add any file you want there, but it's not recommended to add any file that is not related to development since the code there is meant to be ignored by the build process.

For instance there you'll find a `testHelpers.go` file that is used to read the `.env` during tests.

### Building

## REST API

### Authentication

The API accepts two types of tokens and uses the standard "Authorization: Bearer \<token\>" header:

- **sessionToken**: is a JWT, used for dashboard access and generated with the user login credentials. The token contains the userId, IP address and expiration date. It expires in 3 hours and only the IP address used on the token generation is allowed to use it.

- **userApiKey**: is a token meant for M2M communication. The token is a _AES-256-CTR-Encrypted-Base64-Encoded_ string, but only the SHA3-256 hash of the key is stored in the server. The userId is retrieved during key decoding, thus you don't need to provide it. The token never expires, but the user can update it at any time.

### OpenApi // Swagger

To generate the swagger documentation, you must use the following command:

```
swag init -g src/presentation/api/api.go -o src/presentation/api/docs
```

The annotations are in the controller files. The reference file can be found [here](https://github.com/swaggo/swag#attribute).
