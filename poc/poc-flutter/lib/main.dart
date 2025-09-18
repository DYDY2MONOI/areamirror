import 'dart:math';
import 'package:flutter/material.dart';

void main() => runApp(const MyApp());

class MyApp extends StatelessWidget {
  const MyApp({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(home: Home());
  }
}

class Home extends StatefulWidget {
  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  int number = Random().nextInt(100);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Flutter POC")),
      body: Center(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Text("Nombre: $number", style: const TextStyle(fontSize: 24)),
          ElevatedButton(
            onPressed: () => setState(() => number = Random().nextInt(100)),
            child: const Text("Refresh"),
          )
        ]),
      ),
    );
  }
}
