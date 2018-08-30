import React, { Component } from 'react';
import Navigation from './components/navigation';
import Alarm from './components/alarm';

class App extends Component {
  render() {
    return (
      <div>
      <Navigation />
      <Alarm/>
    </div>
    );
  }
}

export default App;
