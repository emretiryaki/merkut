import React, {Component} from 'react';
import {FormControl, FormGroup, ControlLabel, HelpBlock, Checkbox, Radio, Button} from 'react-bootstrap';


class AddAlertForm extends React.Component {
 
    render() {
      return (
        <div className="container-fluid">
        <form>
        <FormGroup controlId="formControlsText">
          <ControlLabel>Alert Name</ControlLabel>
          <FormControl type="text" placeholder="" />
        </FormGroup>
        <FormGroup controlId="formControlsText">
          <ControlLabel>Alert Name</ControlLabel>
          <FormControl type="text" placeholder="" />
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