# -*- coding: utf-8 -*-
"""
Created on Sat Feb 15 11:18:24 2020

@author: Tanmay Thakur
"""
import pandas as pd
import numpy as np
import pickle

from sklearn.preprocessing import StandardScaler


data = pd.read_csv('features.csv')

"""
input_data = data.drop(['Date','Time'], axis = 1).values
"""

input_data = data[['Vrms ph-n AN Avg','Vrms ph-n BN Avg','Vrms ph-n CN Avg']]

T = 10 
D = input_data.shape[1]
N = len(input_data) - T 

Ntrain = len(input_data) * 4//5


scaler = StandardScaler()
scaler.fit(input_data[:Ntrain + T - 1])
input_data = scaler.transform(input_data)


X_train = np.zeros((Ntrain, T, D))
Y_train = np.zeros((Ntrain, D))

for t in range(Ntrain):
  X_train[t, :, :] = input_data[t:t+T]
  Y_train[t] = input_data.loc[t+T]
  
print(X_train.shape, Y_train.shape)

X_test = np.zeros((N - Ntrain, T, D))
Y_test = np.zeros((N - Ntrain, D))

for u in range(N - Ntrain):
  t = u + Ntrain
  X_test[u, :, :] = input_data[t:t+T]
  Y_test[u] = input_data[t+T]

print(X_test.shape, Y_test.shape)
  
data_dump = X_train, X_test, Y_train, Y_test

pickle_out = open("dict.pickle","wb")
pickle.dump(data_dump, pickle_out)
pickle_out.close()

pickle_out = open("scaler.pickle","wb")
pickle.dump(scaler, pickle_out)
pickle_out.close()