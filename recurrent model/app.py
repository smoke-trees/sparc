# -*- coding: utf-8 -*-
"""
Created on Fri Feb 21 23:02:51 2020

@author: Tanmay Thakur
"""
import pickle
import numpy as np
import tensorflow as tf
import matplotlib.pyplot as plt

from sklearn.metrics import mean_squared_error


physical_devices = tf.config.list_physical_devices('GPU')
try:
  tf.config.experimental.set_memory_growth(physical_devices[0], True)
except:
  print("Invalid device or cannot modify virtual devices once initialized.")
  pass

X_train, y_train = pickle.load(open( "dict.pickle", "rb" ))

model = tf.keras.models.load_model("recurrent_model_initial.h5")

validation_target = y_train[3*len(X_train)//4:]
validation_predictions = []
error = []

# index of first validation input
i = 3*len(X_train)//4

while len(validation_predictions) < len(validation_target) - 1:
  p = model.predict(X_train[i].reshape(1, X_train.shape[1], X_train.shape[2]))[0] 
  i += 1

  error.append(mean_squared_error(p,y_train[i]))
  # update the predictions list
  validation_predictions.append(p)
  
plt.plot(error)

pickle.dump(error, open("error.pickle", "wb"))