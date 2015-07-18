/*
    My Implementation is route mapping

    /user/{name}
    /user/profile
    /user/settings
    /user/{name}/image/{id}
    /users
    /users/scores
    /images
    /images/{id}

*/

var RouteTree = {
    98: [
        "user": {
            23: [
                "{name}": {
                    68: [
                        "image": {
                            23: [
                            "{id}": {}
                            ]
                        }
                    ]
                },
                "setting": {
                    // empty
                },
                "profile": {
                    // empty
                }
            ]
        },
        "users": {
            24: ["scores": {}]
        }
    ],
    105: [
        "images": {
            23: ["{id}": {}]
        }
    ]
}
