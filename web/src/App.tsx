import React, { useState } from 'react';
import Attendees from './components/Attendees';
import Seeds from './components/Seeds';
import Conflicts from './components/Conflicts';

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
        >
          Conflicts
        </button>
      </nav>
      
      {activeTab === 'attendees' && <Attendees />}
      {activeTab === 'seeds' && <Seeds />}
      {activeTab === 'conflicts' && <Conflicts />}
    </div>
  );
}

export default App;
