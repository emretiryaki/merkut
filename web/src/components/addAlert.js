import React, {Component} from 'react';
import {FormControl, FormGroup, ControlLabel, HelpBlock, Checkbox, Radio, Button,Form} from 'react-bootstrap';


class AddAlertForm extends React.Component {
 
  constructor() {
    super();
    this.state = {
      name: '',
      shareholders: [{ name: '' }],
    };
  }

  // ...

  handleShareholderNameChange = (idx) => (evt) => {
    const newShareholders = this.state.shareholders.map((shareholder, sidx) => {
      if (idx !== sidx) return shareholder;
      return { ...shareholder, name: evt.target.value };
    });

    this.setState({ shareholders: newShareholders });
  }

  handleSubmit = (evt) => {
    const { name, shareholders } = this.state;
    alert(`Incorporated: ${name} with ${shareholders.length} shareholders`);
  }

  handleAddShareholder = () => {
    this.setState({
      shareholders: this.state.shareholders.concat([{ name: '' }])
    });
  }

  handleRemoveShareholder = (idx) => () => {
    this.setState({
      shareholders: this.state.shareholders.filter((s, sidx) => idx !== sidx)
    });
  }
  
    render() {
      return (
        <div className="container-fluid">
        <Form horizontal>
        <FormGroup controlId="AlertNameTxt">
          <ControlLabel>Alert Name</ControlLabel>
          <FormControl type="text" placeholder="" />
        </FormGroup>
        <FormGroup controlId="ScheduleDrpDown">
        <ControlLabel>Schedule</ControlLabel>
        <FormControl componentClass="select" placeholder="">
          <option value="select"></option>
          <option value="hourly">hourly</option>
          <option value="daily">daily</option>
          <option value="weekly">weekly</option>
          <option value="interval">interval</option>
        </FormControl>
        </FormGroup>
        <FormGroup controlId="WhenTxt">
          <ControlLabel>When</ControlLabel>
          <FormControl type="text" placeholder="" />
        </FormGroup>
        
        <FormGroup controlId="ConditionTypeDrpDown">
        <ControlLabel>Conditions</ControlLabel>
        <FormControl componentClass="select" placeholder="">
          <option value="select"></option>
          <option value="never">never</option>
          <option value="compare">compare</option>
        </FormControl>
        </FormGroup>
        <FormGroup controlId="ConditionOperatorTxt">
          <ControlLabel>Operator</ControlLabel>
          <FormControl type="text" placeholder="" />
        </FormGroup>
        <FormGroup controlId="FieldTxt">
         <ControlLabel>Field</ControlLabel>&nbsp;&nbsp;&nbsp;
         <Button type="button" onClick={this.handleAddShareholder } className="small">+</Button>
         </FormGroup>
         <FormGroup controlId="FieldTxt">
          {this.state.shareholders.map((shareholder, idx) => (
          <div >
            <FormControl
              type="text"
              value={shareholder.name}
              onChange={this.handleShareholderNameChange(idx)}
            />
            
            <Button type="button" onClick={this.handleRemoveShareholder(idx)} className="small">-</Button>
          </div>

          
        ))}
        </FormGroup>
        
        <FormGroup controlId="FieldTxt"></FormGroup>
        <Button type="submit">
          Submit
        </Button>
      </Form>
      </div>
      );
    }
  }
  
export default AddAlertForm;