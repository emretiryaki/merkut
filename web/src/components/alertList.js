import React from 'react';
import API from '../utilities/api'

class Alert extends React.Component{

  state = {
    alerts: [],
    errors: null
  }

  getAlerts() {
    API.get(`alerts`)
    .then(response => {
      this.setState({
        alerts: response.data,
      });
    })
    // If we catch any errors connecting, let's update accordingly
    .catch(error => this.setState({ error, isLoading: false }));
  }

  componentDidMount(){
    this.getAlerts()
  }

  render() {
    const {alerts} =this.state;

    return (
       <div className="container-fluid">
        <table className="table">
        <thead>
          <tr>
            <th>Id</th>
            <th>Name</th>
            <th>State</th>
            <th>Comment</th>
            <th>Last Fired</th>
            <th>Last Triggered</th>
          </tr>
        </thead>
        <tbody>
        {(
            alerts.map(alert => {
              const { id, name, state, comment, lastFired,lastTriggered } = alert;
              return (  <tr key={id}>
                <td>{id}</td>
                <td>{name}</td>
                <td>{state}</td>
                <td>{comment}</td>
                <td>{lastFired}</td>
                <td>{lastTriggered}</td>
                </tr>);
            })
          ) }

        </tbody>
      </table>

      </div>
    )
  }
}


export default Alert;