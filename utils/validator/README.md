# Validator Package

## **How to use ?**

- You need to import validator v2 library `"bitbucket.org/matchmove/platform-go-boilerplate/utils/validator/v2"`

- This is wrapper that we wrote on `github.com/go-playground/validator/v10`

- Purpose is to having everything covered in our format that we use currently, also we are using singleton pattern for that so no need to initialize it anywhere. It will automatically get initialize if its not.

- There is wrapper function which is [Struct()](./v2/validator.go#L55) which will return errors but not in json friendly format.

- if you need error response in json friedly format then use [StructWithFormattedErrors()](./v2/validator.go#L80)

## Which type of validation are supported ?

- There are lot's of validator provided inbuilt by the library please check the [list](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme)

## How to customize error messages ?

- We have define map with having tag name as key, So in the respective key value. You can define your message respectively, [Please check current implementation](./v2/validator.go#L14)

- Currently, we have covered some in general messages that were in v1.

---

## **FAQ's**

### Why v2 for validator ?

- There is already implementation of validator package which has some limitation, and also we need to write function for validation.

- v2 version we have included another package which comes with all the implemetation that are required to validate.

## Why we can't directly replace existing validator package ?

- Because in most of the services we are using that package and it will required too many changes for each services.

- To avoid that problem we have kept v2 seprately so current validation version would't be affected.

## When can we use v2 ?

- If we are implementing new validation at the time we can use this.

