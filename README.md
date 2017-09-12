# Welcome to Boom

A RESTful api-driven content management microservice written in GO ([Go language](http://www.golang.org/)) using the [Revel Web Framework](https://revel.github.io/).

***Alpha . Not ready for public use***

###The Mission: Removing the Headache of Building Effective, Performant, Efficient, Secure & Easy to Extent Content Management APIs
A single, high performance RESTful content management api to reuse, extend, and configure without knowing any GO code. Just a push button install.

###Boom's Core Values
These are our values in building this RESTful content management microservice.
* Ease of Use: I should be able to onboard and work with the app easily without extensive researching.
* Portable / Reusable: I should be able to move the app to any host easily as well as create new instances for future projects.
* Extend / Configure: I can easily extend and configure without modifying core code.
* Performance: I know that my content services app will provide high quality performance for my readers given the constraints of my server.
* Secure: I know that my app will provide security for my content given the constraints of my configuration and server.



## Get Started
### Add Your Content Types:
* Create Your JSON Content Type Files in "con/content/[contentType].json"

### 1. Save New App.Conf File:
* Open file: conf/app.conf.sample
* Check Your Configs
* Save file as app.conf

### 2. Start the Content API:

    revel run github.com/MoonBabyLabs/boom

### 3. Try Out The Content API

    GET http://localhost:9000/api/[version]/[contentType]
    
You should get back an empty json array because you do not have any content. Woohoo.