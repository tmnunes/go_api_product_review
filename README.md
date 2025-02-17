# GoLang API - Product Review 

## Table of Contents
1. [The Challenge](#challenge)
2. [Architecture](#architecture)
3. [Running the App](#running)
4. [Testing de App](#testing)

## The Challenge <a name="challenge"></a>
The goal is to create an Restful API capable to manage products and product reviews. 
The API should be able to calculate the average review for each product.
Language used is Go. 

* API for products with CRUD and get by identifier actions 
    + product information should not return reviews, only average product rating
* API for reviews with CUD actions
* Endpoints to show product reviews
* Service should notify review service when new review is added, modified or deleted
* Each time review is received, it calculates average rating and stores it into persistent storage
* Product reviews and average product ratings should be cached.
