# Welcome to Boom

A RESTful api-driven content management microservice written in GO ([Go language](http://www.golang.org/)) using the [Revel Web Framework](https://revel.github.io/). 

## Get Started
### Add Your Content Types:
* Create Your JSON Content Type Files in "con/content/[contentType].json"

### 1. Save New App.Conf File:
* Open file: conf/app.conf.sample
* Check Your Configs
* Save file as app.conf

### 2. Start the Content API:

    revel run github.com/MoonBabyLabs/boom

### 3. Try Out The API

    GET http://localhost:9000/api/[version]/[contentType]
    
You should get back an empty json array because you do not have any content. Woohoo.