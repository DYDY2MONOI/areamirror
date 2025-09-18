import React, { useEffect, useState } from 'react';
import { View, Text, FlatList } from 'react-native';

export default function App() {
  const [data, setData] = useState([]);

  useEffect(() => {
    fetch('https://jsonplaceholder.typicode.com/posts?_limit=5')
      .then(res => res.json())
      .then(setData);
  }, []);

  return (
    <View style={{ flex:1, padding:20 }}>
      <Text style={{ fontSize:20, marginBottom:10 }}>React Native POC</Text>
      <FlatList
        data={data}
        keyExtractor={item => item.id.toString()}
        renderItem={({item}) => <Text>- {item.title}</Text>}
      />
    </View>
  );
}
