import React, { useState } from 'react';
import Attendees from './components/Attendees';
import Seeds from './components/Seeds';
import Conflicts from './components/Conflicts';
import SetReporter from './components/SetReporter';

function App() {
  const [activeTab, setActiveTab] = useState('attendees');

  return (
    <div>
      <nav style={{ 
        padding: '16px', 
        borderBottom: '1px solid #ccc', 
        backgroundColor: 'white',
        marginBottom: '0'
      }}>
        <button 
          onClick={() => setActiveTab('attendees')}
          className={activeTab === 'attendees' ? 'button' : 'button-secondary'}
          style={{ marginRight: '8px' }}
        >
          Attendees
        </button>
        <button 
          onClick={() => setActiveTab('seeds')}
          className={activeTab === 'seeds' ? 'button' : 'button-secondary'}
          style={{ marginRight: '8px' }}
        >
          Seeds
        </button>
        <button 
          onClick={() => setActiveTab('conflicts')}
          className={activeTab === 'conflicts' ? 'button' : 'button-secondary'}
          style={{ marginRight: '8px' }}
        >
          Conflicts
        </button>
        <button 
          onClick={() => setActiveTab('sets')}
          className={activeTab === 'sets' ? 'button' : 'button-secondary'}
        >
          Sets
        </button>
      </nav>
      
      {activeTab === 'attendees' && <Attendees />}
      {activeTab === 'seeds' && <Seeds />}
      {activeTab === 'conflicts' && <Conflicts />}
      {activeTab === 'sets' && <SetReporter />}
    </div>
  );
}

export default App;
