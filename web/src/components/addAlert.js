import React, {Component} from 'react';
import {FormControl, FormGroup, ControlLabel, HelpBlock, Checkbox, Radio, Button} from 'react-bootstrap';


class AddAlertForm extends React.Component {
 
    render() {
      return (
        <div className="container-fluid">
        <form>
        <FormGroup controlId="formControlAlertName">
          <ControlLabel>Alert Name</ControlLabel>
          <FormControl type="text" placeholder="" />
        </FormGroup>
        <FormGroup controlId="formControlScheduleDrpDown">
        <ControlLabel>Schedule</ControlLabel>
      <FormControl componentClass="select" placeholder="select">
        <option value="select"></option>
        <option value="hourly">hourly</option>
        <option value="daily">daily</option>
        <option value=""></option>
      </FormControl>
        </FormGroup>
        <Button type="submit">
          Submit
        </Button>
      </form>
      </div>
      );
    }
  }
  
export default AddAlertForm;