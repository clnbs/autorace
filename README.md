[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/clnbs/autorace">
    <img src="assets/logo/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Autorace</h3>

  <p align="center">
    Autorace is a personal project of a multiplayer 2D racing game designed to be fully scalable.
    <br />
    <br />
    <a href="https://github.com/clnbs/autorace/issues">Report Bug</a>
    ·
    <a href="https://github.com/clnbs/autorace/issues">Request Feature</a>
  </p>
</p>


<!-- TABLE OF CONTENTS -->
## Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)
* [Acknowledgements](#acknowledgements)


<!-- ABOUT THE PROJECT -->
## About The Project

![Product Name Screen Shot][product-screenshot]

__*This project is still in a prototype stage.*__

Autorace is a personal project of a multiplayer 2D racing game designed to be fully scalable.

Is it over-engineered ? Totally. Is it, at least, a good game? No. This project was built to show and extends my programming skills because I am currently looking for a job. You can find my résumé [here](https://github.com/clnbs/resume).


### Built With

* [Pixel](https://github.com/faiface/pixel) - 2D game library in Go
* [RabbitMQ](https://www.rabbitmq.com/) - Messages broker
* [Go mod](https://blog.golang.org/using-go-modules) - Dependency Management
* [Google UUID](https://github.com/google/uuid) - UUID creation
* [Fluentd](https://www.fluentd.org/) - Logs centralisation
* [Elasticsearch](https://www.elastic.co/elasticsearch/) - Logs storage
* [Kibana](https://www.elastic.co/kibana) - Logs visualization
* [Redis](https://redis.io/) - Cache and short termed storage
* [Hatchful](https://hatchful.shopify.com/) - Logo creation



<!-- GETTING STARTED -->
## Getting Started

DISCLAIMER: this project is still under heavy development. You should not run Autorace in a production environment.

### Prerequisites

In order to compile and run the server stack, you will need :
 - A Debian based Linux operating system (tested on Debian 10)
 - Docker installed. You can find instructions [here](https://docs.docker.com/get-docker/)
 - Git
 
#### Memory usage and configuration
 
In order to compile and run Autorace full stack, you will need 9 GiB of disk space (recommended).
 
 
To make Elasticsearch be able to start in Docker, you have to modify VM heap map allocation bigger:
```
sudo sysctl -w vm.max_map_count=262144
```
 
Alternatively, you can apply this change in `sysctl.conf`:
```
user@host~$ sudo -i
[sudo] password for user: 
root@host~# echo 'vm.max_map_count=262144' >> /etc/sysctl.conf
``` 


### Installation

1) Clone this repository
```
mkdir -p $GOPATH/src/github.com/clnbs
cd $GOPATH/src/github.com/clnbs
git clone https://github.com/clnbs/autorace
```

2) Build server 
```
make server
```

3) Build clients
```
make linux
make windows
```

4) Optional - Clean build artifacts
```
make clean
``` 

You can get help with this Makefile
```
make
OR
make help
```


<!-- USAGE EXAMPLES -->
## Usage

Deploy this project with ease using the Makefile :
- Starting server stack:
```
make run
```

- Clean deployment artifact to terminate demo server stack:
```
make down
``` 


You can get help with this Makefile:
```
make
OR
make help
```



<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/clnbs/autorace/issues) for a list of proposed features (and known issues).

See the [roadmap](ROADMAP.md) for a list of planned feature 

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Colin Bois - <colin.bois@rocketmail.com>

Project Link: [https://github.com/clnbs/autorace](https://github.com/clnbs/autorace)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements

* [StackOverflow](https://stackoverflow.com/)
* [Gustavo Maciel - Gamasutra](https://www.gamasutra.com/blogs/GustavoMaciel/20131229/207833/Generating_Procedural_Racetracks.php)





<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[issues-url]: https://github.com/clnbs/repo/issues
[license-shield]: https://img.shields.io/github/license/clnbs/repo.svg?style=flat-square
[license-url]: https://github.com/clnbs/repo/blob/master/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=flat-square&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/colin-bois-a5b673105/
[product-screenshot]: assets/screenshot/autorace_2020-10-15_19-00-23.png