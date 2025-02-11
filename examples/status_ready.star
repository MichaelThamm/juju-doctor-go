def init():
    jujudoctor.observe('status_ready', on_juju_status)

def on_juju_status(event):
    # add a custom validation
    
    if event.error == True:
        fail("Validation for juju status output failed.")
    else:
        print("Validation for juju status output succeeded.")
    
    