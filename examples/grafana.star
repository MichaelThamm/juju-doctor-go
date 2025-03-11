def init():
    jujudoctor.observe('status', on_status)
    jujudoctor.observe('show_unit', on_show_unit)

def on_status(event):
    validate_status(event.input)

def on_show_unit(event):
    validate_show_unit(event.input)

def validate_status(status):
    print("Ran STATUS probe: {}".format(status.keys()))

def validate_show_unit(show_unit):
    print("Ran SHOW-UNIT probe: {}".format(show_unit.keys()))
