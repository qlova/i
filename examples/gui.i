import @seed {
    style/colors
}

`Fill Red` := {
    color: rgb(100, 0, 0)
}

main
    user = {
        name: "Quentin"
    }
    
    @user
        name = "Bob"
    }

    print(user.name)

    @seed.NewApp("Widgets & Logic")
        expander()
        @clickme $= button("Click me!", `Fill Red`)
            onclick(script
                text = "You clicked me!"
            })

            @button("Click me too!")
                onclick(script
                    clickme.text = "You clicked my child!"
                    clickme.color = colors.blue
                })
            }
        }
        expander()

        @row()
            expander()
            text("Hello "+user.name)
            expander()
        }

        launch()
    }
}